package buildrtc

import (
	"fmt"
	y "github.com/shiguredo/yspata"
	"os"
	"path/filepath"
	"strings"
)

type IOS struct {
	Conf *Config
}

func NewIOS(conf *Config) *IOS {
	return &IOS{conf}
}

func (n *IOS) DistDirDebug() string {
	return y.Join(n.Conf.DistDir, "ios-debug")
}

func (n *IOS) DistDirRelease() string {
	return y.Join(n.Conf.DistDir, "ios-release")
}

func (n *IOS) DistDirCarthage() string {
	return y.Join(n.Conf.DistDir, "ios-carthage")
}

func (n *IOS) ArchiveDir() string {
	return fmt.Sprintf("sora-webrtc-%s-ios", n.Conf.WebRTCVersion())
}

func (n *IOS) ArchiveZip() string {
	return n.ArchiveDir() + ".zip"
}

func (n *IOS) Carthage() string {
	return n.Conf.IOSFramework
}

func (n *IOS) CarthageZip() string {
	return n.Carthage() + ".zip"
}

func (n *IOS) Solutions() string {
	return "solutions = [\n" +
		"  {\n" +
		"    \"url\": \"https://webrtc.googlesource.com/src.git\",\n" +
		"    \"managed\": False,\n" +
		"    \"name\": \"src\",\n" +
		"    \"deps_file\": \"DEPS\",\n" +
		"    \"custom_deps\": {},\n" +
		"  },\n" +
		"]\n" +
		"target_os = [\"ios\", \"mac\"]\n"
}

func (n *IOS) Build() {
	if n.Conf.Debug {
		n.BuildFramework(BUILD_DEBUG)
	}
	if n.Conf.Release {
		n.BuildFramework(BUILD_RELEASE)
	}
}

func (n *IOS) BuildFramework(conf int) {
	name := BuildConfigName(conf)
	fw := n.Conf.IOSFramework

	y.Printf("Build iOS framework for %s...", name)
	wd, _ := os.Getwd()
	os.Chdir(n.Conf.WebRTCSrcDir)

	_, base := filepath.Split(n.Conf.Path)
	base = strings.TrimSuffix(base, ".json")
	dir := y.Join(n.Conf.BuildDir,
		fmt.Sprintf("build-%s", base),
		fmt.Sprintf("ios-framework-%s", name))
	y.Execf("rm -rf %s %s %s",
		y.Join(dir, fw),
		y.Join(dir, n.Conf.IOSDSYM),
		y.Join(dir, "arm64_libs", fw),
		y.Join(dir, "arm_libs", fw),
		y.Join(dir, "x64_libs", fw))
	n.ExecBuildScript(dir, conf)

	os.Chdir(wd)
	n.GenerateBuildInfo(dir + "/" + fw)
}

func (n *IOS) ExecBuildScript(dir string, conf int) {
	name := BuildConfigName(conf)
	args := []string{n.Conf.Python, n.Conf.IOSBuildScript,
		"-o", dir, "--build_config", name, "--arch"}
	if n.Conf.IOSArchArm64 {
		args = append(args, "arm64")
	}
	if n.Conf.IOSArchArm {
		args = append(args, "arm")
	}
	if n.Conf.IOSArchX64 {
		args = append(args, "x64")
	}
	if n.Conf.IOSBitcode {
		args = append(args, "--bitcode")
	}
	if n.Conf.VP9 {
		args = append(args, "--vp9")
	}
	y.Exec("time", args...)
}

func (n *IOS) GenerateBuildInfo(dir string) {
	fmt.Println("Generate build_info.json...")
	file := y.OpenAll(n.Conf.IOSBuildInfo)
	body := fmt.Sprintf("{\n"+
		"    \"webrtc_version\" : \"%s\",\n"+
		"    \"webrtc_revision\" : \"%s\"\n"+
		"}",
		n.Conf.WebRTCBranch, n.Conf.WebRTCRevision)
	_, err := file.WriteString(body)
	y.FailIf(err, "fail write string")
	file.Close()
	y.ExecIg("cp", n.Conf.IOSBuildInfo, dir)
}

func (n *IOS) Archive() {
	bldDir := n.Conf.BuildDir
	distDir := n.Conf.DistDir
	distDirDg := n.DistDirDebug()
	distDirRl := n.DistDirRelease()
	distDirCr := n.DistDirCarthage()

	// clean
	y.Exec("rm", "-rf", distDir, n.ArchiveDir(), n.ArchiveZip())
	y.Exec("mkdir", distDir)
	y.Exec("mkdir", distDirDg)
	y.Exec("mkdir", distDirRl)
	y.Exec("mkdir", distDirCr)

	if n.Conf.Debug {
		var fwDg = y.Join(bldDir, "ios-framework-debug", n.Conf.IOSFramework)
		var dsymDg = y.Join(bldDir, "ios-framework-debug", n.Conf.IOSDSYM)
		y.Exec("cp", "-r", fwDg, distDirDg)
		y.Exec("cp", "-r", dsymDg, distDirDg)
	} else if n.Conf.Release {
		var fwRl = y.Join(bldDir, "ios-framework-release", n.Conf.IOSFramework)
		var dsymRl = y.Join(bldDir, "ios-framework-release", n.Conf.IOSDSYM)
		y.Exec("cp", "-r", fwRl, distDirRl)
		y.Exec("cp", "-r", dsymRl, distDirRl)

		// carthage
		y.Exec("cp", "-r", fwRl, ".")
		y.Exec("zip", "-rq", n.CarthageZip(), n.Carthage())
		y.Exec("rm", "-rf", n.Carthage())
		y.Exec("mv", n.CarthageZip(), distDirDg)
	}

	// archive
	y.Exec("mv", distDir, n.ArchiveDir())
	y.Exec("zip", "-rq", n.ArchiveZip(), n.ArchiveDir())
}

func (n *IOS) Clean() {
	y.Exec("rm", "-rf", n.ArchiveDir(), n.ArchiveZip())
}

func (n *IOS) Reset() {
	// do nothing
}
