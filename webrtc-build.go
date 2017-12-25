// $ go build webrtc-build.go

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var version = "1.0.1"

var fullVersion string

var isMac = false

var isLinux = false

var configFile = "config.json"

var wd, _ = os.Getwd()

var patchDir = filepath.Join(wd, "patch")

var depotToolsURL = "https://chromium.googlesource.com/chromium/tools/depot_tools.git"

var depotToolsDir = filepath.Join(wd, "webrtc/depot_tools")

var gclient = filepath.Join(depotToolsDir, "gclient")

var WebRTCURL = "https://chromium.googlesource.com/external/webrtc"

var WebRTCDir = filepath.Join(wd, "webrtc")

var WebRTCSourceDir = filepath.Join(WebRTCDir, "src")

var gclientConfig = filepath.Join(WebRTCDir, ".gclient")

var buildDir = filepath.Join(WebRTCDir, "build")

var distDir = filepath.Join(WebRTCDir, "dist")

var distDiriOSDebug = filepath.Join(distDir, "ios-debug")

var distDiriOSRelease = filepath.Join(distDir, "ios-release")

var distDiriOSCarthage = filepath.Join(distDir, "ios-carthage")

var distDirAndroidDebug = filepath.Join(distDir, "android-debug")

var distDirAndroidRelease = filepath.Join(distDir, "android-release")

var iOSBuildScript = filepath.Join(WebRTCSourceDir,
	"tools_webrtc/ios/build_ios_libs.py")

var buildInfo = filepath.Join(buildDir, "build_info.json")

var iOSFrameworkName = "WebRTC.framework"

var iOSDsymName = "WebRTC.dSYM"

var iOSStaticName = "librtc_sdk_objc.a"

var iOSHeaderDir = filepath.Join(WebRTCSourceDir,
	"webrtc/sdk/objc/Framework/Headers/WebRTC")

var iOSArchive string

var iOSArchiveZip string

var iOSCarthageFile = iOSFrameworkName

var iOSCarthageFileZip = iOSCarthageFile + ".zip"

var androidBuildScript = filepath.Join(WebRTCSourceDir,
	"tools_webrtc/android/build_aar.py")

var androidArchive string

var androidArchiveZip string

var androidAARName = "libwebrtc.aar"

func containsString(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func FailIfNotExists(filename string) {
	if !Exists(filename) {
		fmt.Printf("Error: File '%s' is not found\n", filename)
		os.Exit(1)
	}
}

// https://github.com/hnakamur/execcommandexample
func RunCommand(name string, arg ...string) (stdout, stderr string, exitCode int, err error) {
	cmd := exec.Command(name, arg...)
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
		fmt.Printf("%s%s\n", header, scanner.Text())
	}
}

func Exec(name string, arg ...string) string {
	fmt.Printf("# %s %s\n", name, strings.Join(arg, " "))
	stdout, _, _, err := RunCommand(name, arg...)
	FailIf(err)
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

func FailIf(err error) bool {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	return true
}

func FailIf2(_ interface{}, err error) {
	FailIf(err)
}

type Config struct {
	WebRTCBranch   string   `json:"webrtc_branch"`
	WebRTCCommit   string   `json:"webrtc_commit"`
	WebRTCRevision string   `json:"webrtc_revision"`
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

func LoadConfig() {
	raw, err := ioutil.ReadFile(configFile)
	FailIf(err)
	json.Unmarshal(raw, &config)

	info := fmt.Sprintf("%s.%s", config.WebRTCBranch, config.WebRTCCommit)
	iOSArchive = fmt.Sprintf("sora-webrtc-%s-ios", info)
	iOSArchiveZip = iOSArchive + ".zip"
	androidArchive = fmt.Sprintf("sora-webrtc-%s-android", info)
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
	fmt.Printf("Checkout the code with release branch M%s (%s)...\n",
		config.WebRTCBranch, config.WebRTCRevision)

	// fetch コマンドの内容を手動で実行する
	// fetch は中断に対応していない (再実行するとエラーになる)
	os.Chdir(WebRTCDir)
	if isMac {
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
	} else if isLinux {
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
				"target_os = [\"android\", \"linux\"]\n")
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
	fmt.Printf("Build iOS framework for %s...\n", config)
	os.Chdir(WebRTCSourceDir)
	buildDir := filepath.Join(buildDir, fmt.Sprintf("ios-framework-%s", config))
	Execf("rm -rf %s %s %s",
		filepath.Join(buildDir, iOSFrameworkName),
		filepath.Join(buildDir, iOSDsymName),
		filepath.Join(buildDir, "arm64_libs", iOSFrameworkName),
		filepath.Join(buildDir, "arm_libs", iOSFrameworkName),
		filepath.Join(buildDir, "x64_libs", iOSFrameworkName))
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
	if file, err := os.OpenFile(buildInfo, os.O_RDWR|os.O_CREATE, 0755); FailIf(err) {
		body := fmt.Sprintf("{\n"+
			"    \"webrtc_version\" : \"%s\",\n"+
			"    \"webrtc_revision\" : \"%s\"\n"+
			"}",
			config.WebRTCBranch, config.WebRTCRevision)
		FailIf2(file.WriteString(body))
		file.Close()
		ExecIgnore("cp", buildInfo, destDir)
	}
}

func BuildiOSStatic(buildConfig string) {
	fmt.Printf("Build iOS static library for %s...\n", buildConfig)
	os.Chdir(WebRTCSourceDir)
	buildDir := filepath.Join(buildDir, fmt.Sprintf("ios-static-%s", buildConfig))
	includeDir := filepath.Join(buildDir, "include")
	Execf("rm -rf %s %s %s %s",
		filepath.Join(buildDir, iOSStaticName),
		includeDir,
		filepath.Join(buildDir, "arm64_libs", iOSStaticName),
		filepath.Join(buildDir, "arm64_libs", "include"))
	ExecOSBuildScript("static_only", buildDir, buildConfig)

	if !Exists(includeDir) {
		FailIf(os.Mkdir(includeDir, 0755))
	}
	headerDirHandle, err1 := os.Open(iOSHeaderDir)
	FailIf(err1)
	infos, err2 := headerDirHandle.Readdir(0)
	FailIf(err2)
	for _, info := range infos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".h") {
			Execf("cp %s %s", filepath.Join(iOSHeaderDir, info.Name()), includeDir)
		}
	}

	os.Chdir(wd)
	GenerateBuildInfo(buildDir)
}

func BuildAndroidLibrary(buildConfig string) {
	fmt.Printf("Build Android library for %s...\n", buildConfig)

	os.Chdir(WebRTCSourceDir)
	buildDir := filepath.Join(buildDir, fmt.Sprintf("android-%s", buildConfig))
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
	var frameworkDebug = filepath.Join(buildDir,
		"ios-framework-debug", iOSFrameworkName)
	var frameworkRelease = filepath.Join(buildDir,
		"ios-framework-release", iOSFrameworkName)
	var dsymDebug = filepath.Join(buildDir,
		"ios-framework-debug", iOSDsymName)
	var dsymRelease = filepath.Join(buildDir,
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
			patchFile := filepath.Join(patchDir, patch.Patch)
			targetFile := filepath.Join(WebRTCSourceDir, patch.Target)
			ApplyPatch(patchFile, targetFile)
		}
	}

	if isMac {
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
	if isMac {
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
		filepath.Join(WebRTCDir, ".gclient"),
		filepath.Join(WebRTCDir, ".gclient_entries"))
}

func Reset() {
	dirs := []string{"webrtc/depot_tools",
		"webrtc/src",
		"webrtc/src/webrtc",
		"webrtc/src/tools_webrtc"}
	for _, dir := range dirs {
		fmt.Printf("Discard changes of %s...\n", dir)
		Exec("git", "-C", dir, "checkout", "--", ".")
	}
}

func PrintHelp() {
	fmt.Println("Usage: build [options] <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  fetch")
	fmt.Println("        Get depot_tools and source files")
	fmt.Println("  build")
	fmt.Println("        Build libraries")
	fmt.Println("  archive")
	fmt.Println("        Archive libraries")
	fmt.Println("  clean")
	fmt.Println("        Remove all built files and discard all changes")
	fmt.Println("  help")
	fmt.Println("        Print this message")
	fmt.Println("  version")
	fmt.Println("        Print version")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}

func main() {
	if runtime.GOOS == "darwin" {
		isMac = true
	} else if runtime.GOOS == "linux" {
		isLinux = true
	} else {
		fmt.Printf("Error: %s OS is not supported\n", runtime.GOOS)
		os.Exit(1)
	}

	flag.Parse()

	if len(os.Args) <= 1 {
		PrintHelp()
		os.Exit(1)
	}

	LoadConfig()
	fullVersion = fmt.Sprintf("%s-%s-%s", version, runtime.GOOS, runtime.GOARCH)

	scheme := BuildScheme{}
	scheme.Debug = containsString(config.BuildConfig, "debug")
	scheme.Release = containsString(config.BuildConfig, "release")
	if isMac {
		scheme.Framework = containsString(config.IOSTargets, "framework")
		scheme.Static = containsString(config.IOSTargets, "static")
	}

	path := os.Getenv("PATH")
	os.Setenv("PATH", depotToolsDir+":"+path)

	subcmd := flag.Arg(0)
	switch subcmd {
	case "fetch":
		GetDepotTools()
		Fetch()

	case "build":
		Build(scheme)

	case "archive":
		Archive()

	case "clean":
		Clean()
		Reset()

	case "help":
		PrintHelp()

	case "version":
		fmt.Println(fullVersion)

	case "selfdist":
		dist := fmt.Sprintf("sora-webrtc-build-%s", fullVersion)
		patchDir := filepath.Join(dist, "patch")
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
		fmt.Printf("Error: Unknown command: %s\n", subcmd)
		os.Exit(1)
	}
}
