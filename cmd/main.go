package main

import (
	"flag"
	"fmt"
	rtc "github.com/shiguredo/sora-webrtc-build"
	y "github.com/shiguredo/yspata"
	"os"
	"path/filepath"
	"runtime"
)

func printHelp() {
	y.PrintLines(
		"Usage: build [options] <command>",
		"",
		"Commands:",
		"  fetch",
		"        Get depot_tools and source files",
		"  build",
		"        Build libraries",
		"  archive",
		"        Archive libraries",
		"  build-clean",
		"        Remove/restore files generated build command",
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

var confFile = "config.json"

var confOpt = flag.String("config", confFile, "configuration file")

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		printHelp()
		y.Fail()
	}

	conf, err := rtc.LoadConfig(*confOpt)
	if err != nil {
		y.Eprintf("cannot load config: %s", err.Error())
		y.Fail()
	}

	path := os.Getenv("PATH")
	os.Setenv("PATH", conf.DepotToolsDir+":"+path)

	var native rtc.Native
	if y.IsMac {
		native = rtc.NewIOS(conf)
	} else if y.IsLinux {
		native = rtc.NewAndroid(conf)
	} else {
		y.Eprintf("%s OS is not supported\n", runtime.GOOS)
		y.Fail()
	}

	bld := rtc.NewBuilder(conf, native)

	subcmd := flag.Arg(0)
	switch subcmd {
	case "fetch":
		bld.GetDepotTools()
		bld.Fetch()

	case "build":
		if !y.Exists(conf.GclientConf) || !y.Exists(conf.GclientEntries) {
			y.Eprintf("%s or %s are not found. Do './webrtc-build fetch'.",
				conf.GclientConf, conf.GclientEntries)
			y.Fail()
		}
		bld.Build()

	case "archive":
		bld.Archive()

	case "clean":
		bld.Clean()
		bld.Reset()

	case "help":
		printHelp()

	case "version":
		y.Printf("webrtc-build %s, library %s",
			rtc.FullVersion, conf.WebRTCVersion())

	case "selfdist":
		dist := fmt.Sprintf("sora-webrtc-build-%s", rtc.FullVersion)

		wd, _ := os.Getwd()
		curPatchDir, _ := filepath.Rel(wd, conf.PatchDir)
		distPatchDir := y.Join(dist, curPatchDir)
		y.Execf("rm -rf %s %s.zip", dist, dist)
		y.Execf("mkdir %s", dist)
		y.Execf("go build -o webrtc-build cmd/main.go")
		y.Execf("cp webrtc-build %s", dist)
		y.Execf("cp config.json %s", dist)
		y.Execf("cp config-ios-dev.json %s", dist)
		os.MkdirAll(distPatchDir, 0755)
		for _, p := range conf.Patches {
			path := y.Join(curPatchDir, p.Patch)
			y.Execf("cp %s %s", path, distPatchDir)
		}
		y.Execf("tar czf %s.tar.gz %s", dist, dist)

	default:
		y.Eprintf("Unknown command: %s", subcmd)
		y.Fail()
	}
}
