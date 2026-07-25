package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jessevdk/go-flags"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/golib/pluginutil"
)

var version string
var commit string

type Opt struct {
	Version bool `short:"v" long:"version" description:"Show version"`
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opt := &Opt{}
	psr := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
	if opt.Version {
		if commit == "" {
			commit = "dev"
		}
		fmt.Printf(
			"%s-%s\n%s/%s, %s, %s\n",
			filepath.Base(os.Args[0]),
			version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			commit)
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	tmpDir := pluginutil.PluginWorkDir()
	u := LinuxUsagePlugin{
		workDir: tmpDir,
	}
	plugin := mp.NewMackerelPlugin(u)
	plugin.Run()
	return 0
}
