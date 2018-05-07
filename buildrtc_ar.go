package buildrtc

import (
	"fmt"
	y "github.com/shiguredo/yspata"
	"io"
	"os"
)

type Android struct {
	Conf *Config
}

func NewAndroid(conf *Config) *Android {
	return &Android{conf}
}

func (n *Android) ArchiveDir() string {
	return fmt.Sprintf("sora-webrtc-%s-android", n.Conf.WebRTCVersion())
}

func (n *Android) ArchiveZip() string {
	return n.ArchiveDir() + ".zip"
}

func (n *Android) Solutions() string {
	return "solutions = [\n" +
		"  {\n" +
		"    \"url\": \"https://webrtc.googlesource.com/src.git\",\n" +
		"    \"managed\": False,\n" +
		"    \"name\": \"src\",\n" +
		"    \"deps_file\": \"DEPS\",\n" +
		"    \"custom_deps\": {},\n" +
		"  },\n" +
		"]\n" +
		"target_os = [\"android\", \"unix\"]\n"
}

func (n *Android) Build() {
	if n.Conf.Debug {
		n.BuildAAR("debug")
	}
	if n.Conf.Release {
		n.BuildAAR("release")
	}
}

func (n *Android) BuildAAR(conf string) {
	y.Printf("Build Android AAR for %s...", conf)

	wd, _ := os.Getwd()
	os.Chdir(n.Conf.WebRTCSrcDir)
	bldDir := y.Join(n.Conf.BuildDir, fmt.Sprintf("android-%s", conf))
	tempDir := y.Join(bldDir, "build")
	libaar := y.Join(bldDir, n.Conf.AndroidAAR)
	y.Execf("mkdir -p %s", bldDir)

	args := []string{n.Conf.Python, n.Conf.AndroidBuildScript,
		"--output", libaar, "--build-dir", tempDir,
		"--build_config", conf, "--arch"}
	if n.Conf.AndroidArchV7A {
		args = append(args, "armeabi-v7a")
	}
	if n.Conf.AndroidArchV8A {
		args = append(args, "arm64-v8a")
	}
	cmd := y.Command("time", args...)
	cmd.OnStdin = func(w io.WriteCloser) {
		io.WriteString(w, "y\n")
	}
	cmd.Run().FailIf("build failed")

	os.Chdir(wd)
}

func (n *Android) Archive() {
	bldDir := n.Conf.BuildDir
	distDir := n.Conf.DistDir
	distDirDg := y.Join(distDir, "android-debug")
	distDirRl := y.Join(distDir, "android-release")

	// clean
	y.Exec("rm", "-rf", distDir, n.ArchiveDir(), n.ArchiveZip())
	y.Exec("mkdir", distDir)
	y.Exec("mkdir", distDirDg)
	y.Exec("mkdir", distDirRl)

	// library
	y.Exec("cp", y.Join(bldDir, "android-debug", n.Conf.AndroidAAR), distDirDg)
	y.Exec("cp", y.Join(bldDir, "android-release", n.Conf.AndroidAAR), distDirRl)

	// archive
	y.Exec("mv", distDir, n.ArchiveDir())
	y.Exec("zip", "-rq", n.ArchiveZip(), n.ArchiveDir())
}

func (n *Android) Clean() {
	y.Exec("rm", "-rf", n.ArchiveDir(), n.ArchiveZip())
}

func (n *Android) Reset() {
	// do nothing
}
