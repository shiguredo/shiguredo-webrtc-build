package buildrtc

import (
	"encoding/json"
	"fmt"
	y "github.com/shiguredo/yspata"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Version = "2.1.0"

var FullVersion = y.FullVersion(Version)

const (
	BUILD_DEBUG = iota
	BUILD_RELEASE
)

func BuildConfigName(conf int) string {
	switch conf {
	case BUILD_DEBUG:
		return "debug"
	case BUILD_RELEASE:
		return "release"
	default:
		panic("invalid")
	}
}

type Config struct {
	Path string

	WebRTCBranch   string `json:"webrtc_branch"`
	WebRTCCommit   string `json:"webrtc_commit"`
	WebRTCRevision string `json:"webrtc_revision"`
	MaintVersion   string `json:"maint_version"`
	Debug          bool   `json:"debug"`
	Release        bool   `json:"release"`
	VP9            bool   `json:"vp9"`

	DepotToolsURL string `json:"_depot_tools_url"`
	WebRTCURL     string `json:"_webrtc_url"`

	Git            string `json:"_git"`
	Python         string `json:"_python"`
	WebRTCDir      string `json:"_webrtc_dir"`
	WebRTCSrcDir   string `json:"_webrtc_src_dir"`
	BuildDir       string `json:"_build_dir"`
	DistDir        string `json:"_dist_dir"`
	PatchDir       string `json:"_patch_dir"`
	DepotToolsDir  string `json:"_depot_tools_dir"`
	Gclient        string `json:"_gclient"`
	GclientConf    string `json:"_gclient_config"`
	GclientEntries string `json:"_gclient_entries"`

	IOSArchArm64   bool   `json:"ios_arch_arm64"`
	IOSArchArm     bool   `json:"ios_arch_arm"`
	IOSArchX64     bool   `json:"ios_arch_x64"`
	IOSTargetFw    bool   `json:"ios_target_framework"`
	IOSTargetSt    bool   `json:"ios_target_static"`
	IOSBitcode     bool   `json:"ios_bitcode"`
	IOSBuildInfo   string `json:"_ios_build_info"`
	IOSBuildScript string `json:"_ios_build_script"`
	IOSFramework   string `json:"_ios_framework"`
	IOSDSYM        string `json:"_ios_dsym"`
	IOSHeaderDir   string `json:"_ios_header_dir"`
	IOSStatic      string `json:"_ios_static"`

	AndroidArchV7A     bool   `json:"android_arch_v7a"`
	AndroidArchV8A     bool   `json:"android_arch_v8a"`
	AndroidAAR         string `json:"_android_aar"`
	AndroidBuildScript string `json:"_android_build_script"`

	ApplyPatch bool    `json:"_apply_patch"`
	Patches    []Patch `json:"_patches"`

	Clean []string `json:"_clean"`
	Reset []string `json:"_reset"`
}

type Patch struct {
	Patch  string `json:"patch"`
	Target string `json:"target"`
}

type Builder struct {
	Conf   *Config
	Native Native
}

type Native interface {
	Solutions() string
	Build()
	Archive()
	Clean()
	Reset()
}

func LoadConfig(path string) (*Config, error) {
	var conf Config
	conf.Path = path
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(raw, &conf)
	setAbsPath(&conf.DepotToolsDir)
	setAbsPath(&conf.WebRTCDir)
	setAbsPath(&conf.WebRTCSrcDir)
	setAbsPath(&conf.BuildDir)
	setAbsPath(&conf.DistDir)
	setAbsPath(&conf.PatchDir)
	conf.IOSHeaderDir = y.Join(conf.WebRTCSrcDir, conf.IOSHeaderDir)
	conf.Gclient = y.Join(conf.DepotToolsDir, conf.Gclient)
	conf.GclientConf = y.Join(conf.WebRTCDir, conf.GclientConf)
	conf.GclientEntries = y.Join(conf.WebRTCDir, conf.GclientEntries)
	return &conf, nil
}

func (c *Config) WebRTCVersion() string {
	return fmt.Sprintf("%s.%s.%s", c.WebRTCBranch, c.WebRTCCommit, c.MaintVersion)
}

func NewBuilder(conf *Config, native Native) *Builder {
	path := os.Getenv("PATH")
	os.Setenv("PATH", conf.DepotToolsDir+":"+path)
	setAbsPath(&conf.WebRTCSrcDir)
	return &Builder{Conf: conf, Native: native}
}

func setAbsPath(path *string) {
	abs, err := filepath.Abs(*path)
	if err != nil {
		panic(fmt.Sprintf("cannot get absolute path for %s", path))
	}
	*path = abs
}

func (b *Builder) GetDepotTools() {
	if !y.Exists(b.Conf.DepotToolsDir) {
		fmt.Println("Get depot_tools...")
		y.Exec(b.Conf.Git, "clone", b.Conf.DepotToolsURL, b.Conf.DepotToolsDir)
	} else {
		fmt.Println("Update depot_tools...")
		y.Exec(b.Conf.Git, "-C", b.Conf.DepotToolsDir, "pull")
	}
}

func (b *Builder) Fetch() {
	y.Printf("Checkout the code with release branch M%s (%s)...",
		b.Conf.WebRTCBranch, b.Conf.WebRTCRevision)

	wd, _ := os.Getwd()

	os.Chdir(b.Conf.WebRTCDir)
	y.Exec(b.Conf.Gclient, "config", "--spec", b.Native.Solutions())
	y.Exec(b.Conf.Gclient, "sync", "--nohooks", "--with_branch_heads", "-v", "-R")

	os.Chdir(b.Conf.WebRTCSrcDir)
	y.Exec(b.Conf.Git, "submodule", "foreach", "'git config -f $toplevel/.git/config submodule.$name.ignore all'")
	y.Exec(b.Conf.Git, "config", "diff.ignoreSubmodules", "all")

	y.Exec(b.Conf.Git, "fetch", "origin")
	y.Exec(b.Conf.Git, "checkout", "-B",
		fmt.Sprintf("M%s", b.Conf.WebRTCBranch),
		fmt.Sprintf("remotes/branch-heads/%s", b.Conf.WebRTCBranch))
	y.Exec(b.Conf.Git, "checkout", b.Conf.WebRTCRevision)
	syncCmd2 := y.Command(b.Conf.Gclient, "sync", "--with_branch_heads", "-v", "-R")
	syncCmd2.OnStdin = func(w io.WriteCloser) {
		io.WriteString(w, "y\n")
	}
	syncCmd2.Run().FailIf("build failed, gclient sync")

	y.Exec(b.Conf.Gclient, "runhooks", "-v")

	os.Chdir(wd)
}

func (b *Builder) ApplyPatch(patch string, target string) {
	patch2 := y.Join(b.Conf.PatchDir, patch)
	target2 := y.Join(b.Conf.WebRTCSrcDir, target)
	y.FailIfNotExists(patch2)
	y.FailIfNotExists(target2)
	y.ExecIg("patch", "-buN", target2, patch2)
}

func (b *Builder) Build() {
	if b.Conf.ApplyPatch {
		fmt.Println("Apply patches...")
		for _, p := range b.Conf.Patches {
			b.ApplyPatch(p.Patch, p.Target)
		}
	}

	b.Native.Build()
}

func (b *Builder) Archive() {
	b.Native.Archive()
}

func (b *Builder) Clean() {
	files := []string{b.Conf.DepotToolsDir, b.Conf.Gclient, b.Conf.GclientEntries}
	for _, f := range b.Conf.Clean {
		files = append(files, f)
	}
	for _, f := range files {
		y.Exec("rm", "-rf", f)
	}
	b.Native.Clean()

}

func (b *Builder) BuildClean() {
	files := []string{b.Conf.BuildDir}
	for _, f := range files {
		y.Exec("rm", "-rf", f)
	}
	b.Reset()
}

func (b *Builder) Reset() {
	files := []string{b.Conf.DepotToolsDir, b.Conf.WebRTCSrcDir}
	for _, f := range b.Conf.Reset {
		files = append(files, f)
	}
	for _, f := range files {
		if y.Exists(f) {
			y.Printf("Discard changes of %s...", f)
			y.Exec(b.Conf.Git, "-C", f, "checkout", "--", ".")
		} else {
			y.Printf("Discard changes of %s (not found, ignore)", f)
		}
	}
	b.Native.Reset()
}
