package main

import (
	"os"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/golib/pluginutil"
	"github.com/monitoring-forge/flagrun"
)

var version string

type Opt struct {
	Version bool `short:"v" long:"version" description:"Show version"`
	workDir string
}

func main() {
	opt := &Opt{
		workDir: pluginutil.PluginWorkDir(),
	}
	os.Exit(flagrun.Ship(opt, flagrun.Version(version)))
}

func (opt *Opt) Run(_ []string) {
	plugin := mp.NewMackerelPlugin(opt)
	plugin.Run()
}
