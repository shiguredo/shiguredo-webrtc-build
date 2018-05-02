package main

import (
	"flag"
	"fmt"
	rtc "github.com/shiguredo/sora-webrtc-build"
	y "github.com/shiguredo/yspata"
	"os"
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
			rtc.FullVersion, conf.WebRTCVersion)

	case "selfdist":
		dist := fmt.Sprintf("sora-webrtc-build-%s", rtc.FullVersion)
		patchDir := y.Join(dist, "patch")
		y.Execf("rm -rf %s %s.zip", dist, dist)
		y.Execf("mkdir %s", dist)
		y.Execf("go build webrtc-build.go")
		y.Execf("cp webrtc-build %s", dist)
		y.Execf("cp config.json %s", dist)
		os.MkdirAll(patchDir, 0755)
		y.Execf("cp patch/webrtc_tools_BUILD.gn.diff %s", patchDir)
		y.Execf("cp patch/webrtc_webrtc.gni.diff %s", patchDir)
		y.Execf("cp patch/build_ios_libs.py.diff %s", patchDir)
		y.Execf("cp patch/build_aar.py.diff %s", patchDir)
		y.Execf("tar czf %s.tar.gz %s", dist, dist)

	default:
		y.Eprintf("Unknown command: %s", subcmd)
		y.Fail()
	}
}
