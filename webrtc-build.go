// $ go build webrtc-build.go

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/shiguredo/yspata"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var version = "1.0.1"

var fullVersion = FullVersion(version)

var defaultConfigFile = "config.json"

var wd, _ = os.Getwd()

var patchDir = Join(wd, "patch")

var depotToolsURL = "https://chromium.googlesource.com/chromium/tools/depot_tools.git"

var depotToolsDir = Join(wd, "webrtc/depot_tools")

var gclient = Join(depotToolsDir, "gclient")

var WebRTCURL = "https://chromium.googlesource.com/external/webrtc"

var WebRTCDir = Join(wd, "webrtc")

var WebRTCSourceDir = Join(WebRTCDir, "src")

var gclientConfig = Join(WebRTCDir, ".gclient")

var gclientEntries = Join(WebRTCDir, ".gclient_entries")

var buildDir = Join(WebRTCDir, "build")

var distDir = Join(WebRTCDir, "dist")

var distDiriOSDebug = Join(distDir, "ios-debug")

var distDiriOSRelease = Join(distDir, "ios-release")

var distDiriOSCarthage = Join(distDir, "ios-carthage")

var distDirAndroidDebug = Join(distDir, "android-debug")

var distDirAndroidRelease = Join(distDir, "android-release")

var iOSBuildScript = Join(WebRTCSourceDir,
	"tools_webrtc/ios/build_ios_libs.py")

var buildInfo = Join(buildDir, "build_info.json")

var iOSFrameworkName = "WebRTC.framework"

var iOSDsymName = "WebRTC.dSYM"

var iOSStaticName = "librtc_sdk_objc.a"

var iOSHeaderDir = Join(WebRTCSourceDir,
	"webrtc/sdk/objc/Framework/Headers/WebRTC")

var iOSArchive string

var iOSArchiveZip string

var iOSCarthageFile = iOSFrameworkName

var iOSCarthageFileZip = iOSCarthageFile + ".zip"

var androidBuildScript = Join(WebRTCSourceDir,
	"tools_webrtc/android/build_aar.py")

var androidArchive string

var androidArchiveZip string

var androidAARName = "libwebrtc.aar"

// https://github.com/hnakamur/execcommandexample
func RunCommand(name string, arg ...string) (stdout, stderr string, exitCode int, err error) {
	cmd := exec.Command(name, arg...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	io.WriteString(stdin, "y\n")
	stdin.Close()

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	var bufout, buferr bytes.Buffer
	outReader2 := io.TeeReader(outReader, &bufout)
	errReader2 := io.TeeReader(errReader, &buferr)

	if err = cmd.Start(); err != nil {
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() { PrintOutputWithHeader("stdout:", outReader2); wg.Done() }()
	go func() { PrintOutputWithHeader("stderr:", errReader2); wg.Done() }()
	wg.Wait()
	err = cmd.Wait()

	stdout = bufout.String()
	stderr = buferr.String()

	if err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				err = nil
				exitCode = s.ExitStatus()
			}
		}
	}
	return
}

func PrintOutputWithHeader(header string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		Printf("%s%s", header, scanner.Text())
	}
}

func Exec(name string, arg ...string) string {
	fmt.Printf("# %s %s\n", name, strings.Join(arg, " "))
	stdout, _, _, err := RunCommand(name, arg...)
	FailIf(err, "exec failed")
	return stdout
}

func Execf(format string, arg ...interface{}) {
	cmd := fmt.Sprintf(format, arg...)
	args := strings.Split(cmd, " ")
	Exec(args[0], args[1:]...)
}

func ExecIgnore(name string, arg ...string) {
	fmt.Printf("# %s %s\n", name, strings.Join(arg, " "))
	RunCommand(name, arg...)
}

func ExecIgnoref(format string, arg ...interface{}) {
	cmd := fmt.Sprintf(format, arg...)
	args := strings.Split(cmd, " ")
	ExecIgnore(args[0], args[1:]...)
}

type Config struct {
	WebRTCBranch   string   `json:"webrtc_branch"`
	WebRTCCommit   string   `json:"webrtc_commit"`
	WebRTCRevision string   `json:"webrtc_revision"`
	MaintVersion   string   `json:"maint_version"`
	Python         string   `json:"python"`
	IOSArch        []string `json:"ios_arch"`
	IOSTargets     []string `json:"ios_targets"`
	IOSBitcode     bool     `json:"ios_bitcode"`
	AndroidArch    []string `json:"android_arch"`
	BuildConfig    []string `json:"build_config"`
	VP9            bool     `json:"vp9"`
	ApplyPatch     bool     `json:"apply_patch"`
	Patches        []Patch  `json:"patches"`
}

type Patch struct {
	Patch  string `json:"patch"`
	Target string `json:"target"`
}

var config Config
var webRTCLibVersion string

func LoadConfig() {
	raw, err := ioutil.ReadFile(*configOpt)
	FailIf(err, "cannot read config file")
	json.Unmarshal(raw, &config)

	webRTCLibVersion = fmt.Sprintf("%s.%s.%s", config.WebRTCBranch, config.WebRTCCommit, config.MaintVersion)
	iOSArchive = fmt.Sprintf("sora-webrtc-%s-ios", webRTCLibVersion)
	iOSArchiveZip = iOSArchive + ".zip"
	androidArchive = fmt.Sprintf("sora-webrtc-%s-android", webRTCLibVersion)
	androidArchiveZip = androidArchive + ".zip"
}

func GetDepotTools() {
	if !Exists(depotToolsDir) {
		fmt.Println("Get depot_tools...")
		Exec("git", "clone", depotToolsURL, depotToolsDir)
	} else {
		fmt.Println("Update depot_tools...")
		Exec("git", "-C", depotToolsDir, "pull")
	}
}

func Fetch() {
	Printf("Checkout the code with release branch M%s (%s)...",
		config.WebRTCBranch, config.WebRTCRevision)

	// fetch コマンドの内容を手動で実行する
	// fetch は中断に対応していない (再実行するとエラーになる)
	os.Chdir(WebRTCDir)
	if IsMac {
		Exec(gclient, "config", "--spec",
			"solutions = [\n"+
				"  {\n"+
				"    \"url\": \"https://webrtc.googlesource.com/src.git\",\n"+
				"    \"managed\": False,\n"+
				"    \"name\": \"src\",\n"+
				"    \"deps_file\": \"DEPS\",\n"+
				"    \"custom_deps\": {},\n"+
				"  },\n"+
				"]\n"+
				"target_os = [\"ios\", \"mac\"]\n")
	} else if IsLinux {
		Exec(gclient, "config", "--spec",
			"solutions = [\n"+
				"  {\n"+
				"    \"url\": \"https://webrtc.googlesource.com/src.git\",\n"+
				"    \"managed\": False,\n"+
				"    \"name\": \"src\",\n"+
				"    \"deps_file\": \"DEPS\",\n"+
				"    \"custom_deps\": {},\n"+
				"  },\n"+
				"]\n"+
				"target_os = [\"android\", \"unix\"]\n")
	} else {
		panic("unsupported OS")
	}

	Exec(gclient, "sync", "--nohooks", "--with_branch_heads", "-v", "-R")
	Exec("git", "submodule", "foreach", "'git config -f $toplevel/.git/config submodule.$name.ignore all'")
	Exec("git", "config", "--add", "remote.origin.fetch", "'+refs/tags/*:refs/tags/*'")
	Exec("git", "config", "diff.ignoreSubmodules", "all")

	// end fetch

	os.Chdir(WebRTCSourceDir)
	Exec("git", "fetch", "origin")
	Exec("git", "checkout", "-B",
		fmt.Sprintf("M%s", config.WebRTCBranch),
		fmt.Sprintf("refs/remotes/branch-heads/%s", config.WebRTCBranch))
	Exec("git", "checkout", config.WebRTCRevision)
	Exec(gclient, "sync", "--with_branch_heads", "-v", "-R")
	Exec(gclient, "runhooks", "-v")
	os.Chdir(wd)
}

func ApplyPatch(patch string, target string) {
	FailIfNotExists(patch)
	FailIfNotExists(target)
	ExecIgnore("patch", "-buN", target, patch)
}

func BuildiOSFramework(config string) {
	Printf("Build iOS framework for %s...", config)
	os.Chdir(WebRTCSourceDir)
	buildDir := Join(buildDir, fmt.Sprintf("ios-framework-%s", config))
	Execf("rm -rf %s %s %s",
		Join(buildDir, iOSFrameworkName),
		Join(buildDir, iOSDsymName),
		Join(buildDir, "arm64_libs", iOSFrameworkName),
		Join(buildDir, "arm_libs", iOSFrameworkName),
		Join(buildDir, "x64_libs", iOSFrameworkName))
	ExecOSBuildScript("framework", buildDir, config)
	os.Chdir(wd)
	GenerateBuildInfo(buildDir + "/" + iOSFrameworkName)
}

func ExecOSBuildScript(ty string, buildDir, buildConfig string) {
	args := []string{config.Python, iOSBuildScript, "-o", buildDir,
		"-b", ty, "--build_config", buildConfig, "--arch"}
	args = append(args, config.IOSArch...)
	if config.IOSBitcode {
		args = append(args, "--bitcode")
	}
	if config.VP9 {
		args = append(args, "--vp9")
	}
	Exec("time", args...)
}

func GenerateBuildInfo(destDir string) {
	fmt.Println("Generate build_info.json...")
	file := OpenFile(buildInfo, os.O_RDWR|os.O_CREATE, 0755)
	body := fmt.Sprintf("{\n"+
		"    \"webrtc_version\" : \"%s\",\n"+
		"    \"webrtc_revision\" : \"%s\"\n"+
		"}",
		config.WebRTCBranch, config.WebRTCRevision)
	_, err := file.WriteString(body)
	FailIf(err, "fail write string")
	file.Close()
	ExecIgnore("cp", buildInfo, destDir)
}

func BuildiOSStatic(buildConfig string) {
	Printf("Build iOS static library for %s...", buildConfig)
	os.Chdir(WebRTCSourceDir)
	buildDir := Join(buildDir, fmt.Sprintf("ios-static-%s", buildConfig))
	includeDir := Join(buildDir, "include")
	Execf("rm -rf %s %s %s %s",
		Join(buildDir, iOSStaticName),
		includeDir,
		Join(buildDir, "arm64_libs", iOSStaticName),
		Join(buildDir, "arm64_libs", "include"))
	ExecOSBuildScript("static_only", buildDir, buildConfig)

	if !Exists(includeDir) {
		FailIf(os.Mkdir(includeDir, 0755), "cannot mkdir")
	}
	headerDirHandle, err1 := os.Open(iOSHeaderDir)
	FailIf(err1, "cannot open %s", iOSHeaderDir)
	infos, err2 := headerDirHandle.Readdir(0)
	FailIf(err2, "cannot read")
	for _, info := range infos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".h") {
			Execf("cp %s %s", Join(iOSHeaderDir, info.Name()), includeDir)
		}
	}

	os.Chdir(wd)
	GenerateBuildInfo(buildDir)
}

func BuildAndroidLibrary(buildConfig string) {
	Printf("Build Android library for ...", buildConfig)

	os.Chdir(WebRTCSourceDir)
	buildDir := Join(buildDir, fmt.Sprintf("android-%s", buildConfig))
	tempDir := buildDir + "/build"
	libaar := buildDir + "/" + androidAARName
	Execf("mkdir -p %s", buildDir)
	args := []string{config.Python, androidBuildScript,
		"--output", libaar, "--build-dir", tempDir,
		"--build_config", buildConfig, "--arch"}
	args = append(args, config.AndroidArch...)
	Exec("time", args...)
	os.Chdir(wd)
}

func ArchiveiOSProducts() {
	// clean
	Exec("rm", "-rf", distDir, iOSArchive, iOSArchiveZip)
	Exec("mkdir", distDir)
	Exec("mkdir", distDiriOSDebug)
	Exec("mkdir", distDiriOSRelease)
	Exec("mkdir", distDiriOSCarthage)

	// framework
	var frameworkDebug = Join(buildDir,
		"ios-framework-debug", iOSFrameworkName)
	var frameworkRelease = Join(buildDir,
		"ios-framework-release", iOSFrameworkName)
	var dsymDebug = Join(buildDir,
		"ios-framework-debug", iOSDsymName)
	var dsymRelease = Join(buildDir,
		"ios-framework-release", iOSDsymName)
	Exec("cp", "-r", frameworkDebug, distDiriOSDebug)
	Exec("cp", "-r", frameworkRelease, distDiriOSRelease)
	Exec("cp", "-r", dsymDebug, distDiriOSDebug)
	Exec("cp", "-r", dsymRelease, distDiriOSRelease)

	// carthage
	Exec("cp", "-r", frameworkRelease, ".")
	Exec("zip", "-rq", iOSCarthageFileZip, iOSCarthageFile)
	Exec("rm", "-rf", iOSCarthageFile)
	Exec("mv", iOSCarthageFileZip, distDiriOSCarthage)

	// static
	Exec("cp", "-r", buildDir+"/ios-static-debug/arm64_libs/librtc_sdk_objc.a",
		distDiriOSDebug)
	Exec("cp", "-r", buildDir+"ios-static-debug/include", distDiriOSDebug)
	Exec("cp", "-r", buildDir+"/ios-static-release/arm64_libs/librtc_sdk_objc.a",
		distDiriOSRelease)
	Exec("cp", "-r", buildDir+"/ios-static-release/include", distDiriOSRelease)

	// archive
	Exec("mv", distDir, iOSArchive)
	Exec("zip", "-rq", iOSArchiveZip, iOSArchive)
}

func ArchiveAndroidProducts() {
	// clean
	Exec("rm", "-rf", distDir, androidArchive, androidArchiveZip)
	Exec("mkdir", distDir)
	Exec("mkdir", distDirAndroidDebug)
	Exec("mkdir", distDirAndroidRelease)

	// library
	Exec("cp", buildDir+"/android-debug/"+androidAARName, distDirAndroidDebug)
	Exec("cp", buildDir+"/android-release/"+androidAARName, distDirAndroidRelease)

	// archive
	Exec("mv", distDir, androidArchive)
	Exec("zip", "-rq", androidArchiveZip, androidArchive)
}

type BuildScheme struct {
	Debug, Release, Framework, Static bool
}

func Build(scheme BuildScheme) {
	if config.ApplyPatch {
		fmt.Println("Apply patches...")
		for _, patch := range config.Patches {
			patchFile := Join(patchDir, patch.Patch)
			targetFile := Join(WebRTCSourceDir, patch.Target)
			ApplyPatch(patchFile, targetFile)
		}
	}

	if IsMac {
		if scheme.Framework {
			if scheme.Debug {
				BuildiOSFramework("debug")
			}
			if scheme.Release {
				BuildiOSFramework("release")
			}
		}
		if scheme.Static {
			if scheme.Debug {
				BuildiOSStatic("debug")
			}
			if scheme.Release {
				BuildiOSStatic("release")
			}
		}
	} else {
		if scheme.Debug {
			BuildAndroidLibrary("debug")
		}
		if scheme.Release {
			BuildAndroidLibrary("release")
		}
	}
}

func Archive() {
	if IsMac {
		ArchiveiOSProducts()
	} else {
		ArchiveAndroidProducts()
	}
}

func Clean() {
	Exec("rm", "-rf", buildDir,
		iOSArchive, iOSArchiveZip, androidArchive, androidArchiveZip,
		"webrtc/src/testing/gmock",
		"webrtc/src/testing/gtest",
		Join(WebRTCDir, ".gclient"),
		Join(WebRTCDir, ".gclient_entries"))
}

func Reset() {
	dirs := []string{"webrtc/depot_tools",
		"webrtc/src",
		"webrtc/src/webrtc",
		"webrtc/src/tools_webrtc"}
	for _, dir := range dirs {
		Printf("Discard changes of %s...", dir)
		Exec("git", "-C", dir, "checkout", "--", ".")
	}
}

func PrintHelp() {
	PrintLines(
		"Usage: build [options] <command>",
		"",
		"Commands:",
		"  fetch",
		"        Get depot_tools and source files",
		"  build",
		"        Build libraries",
		"  archive",
		"        Archive libraries",
		"  clean",
		"        Remove all built files and discard all changes",
		"  help",
		"        Print this message",
		"  version",
		"        Print version",
		"",
		"Options:")
	flag.PrintDefaults()
}

var configOpt = flag.String("config", defaultConfigFile, "configuration file")

func main() {
	if !(IsMac || IsLinux) {
		Eprintf("%s OS is not supported\n", runtime.GOOS)
		os.Exit(1)
	}

	flag.Parse()

	if flag.NArg() < 1 {
		PrintHelp()
		os.Exit(1)
	}

	LoadConfig()

	scheme := BuildScheme{}
	scheme.Debug = Contains(config.BuildConfig, "debug")
	scheme.Release = Contains(config.BuildConfig, "release")
	if IsMac {
		scheme.Framework = Contains(config.IOSTargets, "framework")
		scheme.Static = Contains(config.IOSTargets, "static")
	}

	path := os.Getenv("PATH")
	os.Setenv("PATH", depotToolsDir+":"+path)

	subcmd := flag.Arg(0)
	switch subcmd {
	case "fetch":
		GetDepotTools()
		Fetch()

	case "build":
		if !Exists(gclientConfig) || !Exists(gclientEntries) {
			Eprintf("webrtc/.gclient or webrtc/.gclient_entries are not found. Do './webrtc-build fetch'.")
			os.Exit(1)
		}
		Build(scheme)

	case "archive":
		Archive()

	case "clean":
		Clean()
		Reset()

	case "help":
		PrintHelp()

	case "version":
		Printf("webrtc-build %s, library %s", fullVersion, webRTCLibVersion)

	case "selfdist":
		dist := fmt.Sprintf("sora-webrtc-build-%s", fullVersion)
		patchDir := Join(dist, "patch")
		Execf("rm -rf %s %s.zip", dist, dist)
		Execf("mkdir %s", dist)
		Execf("go build webrtc-build.go")
		Execf("cp webrtc-build %s", dist)
		Execf("cp config.json %s", dist)
		os.MkdirAll(patchDir, 0755)
		Execf("cp patch/webrtc_tools_BUILD.gn.diff %s", patchDir)
		Execf("cp patch/webrtc_webrtc.gni.diff %s", patchDir)
		Execf("cp patch/build_ios_libs.py.diff %s", patchDir)
		Execf("cp patch/build_aar.py.diff %s", patchDir)
		Execf("tar czf %s.tar.gz %s", dist, dist)

	default:
		Eprintf("Unknown command: %s", subcmd)
		os.Exit(1)
	}
}
