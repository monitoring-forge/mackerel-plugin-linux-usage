package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jessevdk/go-flags"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/golib/pluginutil"
	"github.com/prometheus/procfs"
)

var version string
var commit string

type Opt struct {
	Version bool `short:"v" long:"version" description:"Show version"`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func generateTempfilePath() string {
	tmpDir := pluginutil.PluginWorkDir()
	curUser, _ := user.Current()
	uid := "0"
	if curUser != nil {
		uid = curUser.Uid
	}
	path := filepath.Join(tmpDir, fmt.Sprintf("mackerel-plugin-linux-usage-%s", uid))
	return path
}

type stats struct {
	CPUStat procfs.CPUStat `json:"cpustat"`
	Time    int64          `json:"time"`
}

func writeStats(filename string, st procfs.CPUStat) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	n := time.Now().Unix()
	jb, err := json.Marshal(stats{st, n})
	if err != nil {
		return err
	}
	_, err = file.Write(jb)
	return err
}

func readStats(filename string) (int64, procfs.CPUStat, error) {
	st := stats{}
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, procfs.CPUStat{}, err
	}
	err = json.Unmarshal(d, &st)
	if err != nil {
		return 0, procfs.CPUStat{}, err
	}
	return st.Time, st.CPUStat, nil
}

type LinuxUsagePlugin struct{}

func (u LinuxUsagePlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"linux-usage.cpu": {
			Label: "Linux CPU usage max 100%",
			Unit:  mp.UnitPercentage,
			Metrics: []mp.Metrics{
				{Name: "guest_nice", Label: "guest_nice", Diff: false, Stacked: true},
				{Name: "guest", Label: "guest", Diff: false, Stacked: true},
				{Name: "steal", Label: "steal", Diff: false, Stacked: true},
				{Name: "softirq", Label: "softirq", Diff: false, Stacked: true},
				{Name: "irq", Label: "irq", Diff: false, Stacked: true},
				{Name: "iowait", Label: "ioWait", Diff: false, Stacked: true},
				{Name: "idle", Label: "idle", Diff: false, Stacked: true},
				{Name: "system", Label: "system", Diff: false, Stacked: true},
				{Name: "nice", Label: "nice", Diff: false, Stacked: true},
				{Name: "user", Label: "user", Diff: false, Stacked: true},
			},
		},
		"linux-usage.loadavg": {
			Label: "Linux CPU load average per CPU",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "loadavg1", Label: "loadavg1", Diff: false, Stacked: false},
				{Name: "loadavg5", Label: "loadavg5", Diff: false, Stacked: false},
				{Name: "loadavg15", Label: "loadavg15", Diff: false, Stacked: false},
			},
		},
		"linux-usage.process": {
			Label: "Linux CPU number of processes",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "all", Label: "all", Diff: false, Stacked: false},
				{Name: "running", Label: "running", Diff: false, Stacked: false},
			},
		},
		"linux-usage.tcp-opens": {
			Label: "Linux CPU TCP Opens",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "active", Label: "ActiveOpens", Diff: true, Stacked: false},
				{Name: "passive", Label: "PassiveOpens", Diff: true, Stacked: false},
			},
		},
		"linux-usage.tcp-listen": {
			Label: "Linux CPU TCP Listen",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "overflows", Label: "ListenOverflows", Diff: true, Stacked: false},
				{Name: "drops", Label: "ListenDrops", Diff: true, Stacked: false},
			},
		},
	}
}

func (u LinuxUsagePlugin) FetchMetrics() (map[string]float64, error) {
	res := map[string]float64{}
	pfs, err := procfs.NewDefaultFS()
	if err != nil {
		return res, err
	}
	st, err := pfs.Stat()
	if err != nil {
		return res, err
	}
	loadavg, err := pfs.LoadAvg()
	if err != nil {
		return res, err
	}
	procs, err := pfs.AllProcs()
	if err != nil {
		return res, err
	}
	selffs, err := pfs.Self()
	if err != nil {
		return res, err
	}
	psnmp, err := selffs.Snmp()
	if err != nil {
		return res, err
	}
	pnetstat, err := selffs.Netstat()
	if err != nil {
		return res, err
	}

	totalProcs := float64(0)
	procRunning := float64(0)
	for _, proc := range procs {
		ps, err := proc.Stat()
		if err != nil {
			continue
		}
		if ps.State == "R" {
			procRunning++
		}
		totalProcs++
	}

	cur := st.CPUTotal
	cores := float64(len(st.CPU))

	res["loadavg1"] = loadavg.Load1 / cores
	res["loadavg5"] = loadavg.Load5 / cores
	res["loadavg15"] = loadavg.Load15 / cores
	res["all"] = totalProcs
	res["running"] = procRunning
	res["active"] = *psnmp.Tcp.ActiveOpens
	res["passive"] = *psnmp.Tcp.PassiveOpens
	res["overflows"] = *pnetstat.ListenOverflows
	res["drops"] = *pnetstat.ListenDrops

	path := generateTempfilePath()

	if !fileExists(path) {
		err = writeStats(path, cur)
		if err != nil {
			return res, err
		}
		return res, nil
	}

	defer func() {
		err := writeStats(path, cur)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}()

	t, prev, err := readStats(path)
	if err != nil {
		return res, err
	}
	if t == 0 {
		return res, fmt.Errorf("failed to get previous time")
	}
	n := time.Now().Unix()
	timeDiff := float64(n - t)
	if timeDiff > 600 {
		return res, fmt.Errorf("too long duration")
	}

	var total float64
	// User
	gapUser := float64(cur.User) - float64(prev.User)
	if gapUser < 0 {
		gapUser = 0
	}
	total += gapUser
	// Nice
	gapNice := float64(cur.Nice) - float64(prev.Nice)
	if gapNice < 0 {
		gapNice = 0
	}
	total += gapNice
	// System
	gapSystem := float64(cur.System) - float64(prev.System)
	if gapSystem < 0 {
		gapSystem = 0
	}
	total += gapSystem
	// Idle
	gapIdle := float64(cur.Idle) - float64(prev.Idle)
	if gapIdle < 0 {
		gapIdle = 0
	}
	total += gapIdle
	// Iowait
	gapIowait := float64(cur.Iowait) - float64(prev.Iowait)
	if gapIowait < 0 {
		gapIowait = 0
	}
	total += gapIowait
	// Irq
	gapIRQ := float64(cur.IRQ) - float64(prev.IRQ)
	if gapIRQ < 0 {
		gapIRQ = 0
	}
	total += gapIRQ
	// SoftIRQ
	gapSoftIRQ := float64(cur.SoftIRQ) - float64(prev.SoftIRQ)
	if gapSoftIRQ < 0 {
		gapSoftIRQ = 0
	}
	total += gapSoftIRQ
	// Steal
	gapSteal := float64(cur.Steal) - float64(prev.Steal)
	if gapSteal < 0 {
		gapSteal = 0
	}
	total += gapSteal
	// Guest
	gapGuest := float64(cur.Guest) - float64(prev.Guest)
	if gapGuest < 0 {
		gapGuest = 0
	}
	total += gapGuest
	// GuestNice
	gapGuestNice := float64(cur.GuestNice) - float64(prev.GuestNice)
	if gapGuestNice < 0 {
		gapGuestNice = 0
	}
	total += gapGuestNice

	// User includes Guest
	gapUser -= gapGuest
	total -= gapGuest
	// Nice includes GuestNice
	gapNice -= gapGuestNice
	total -= gapGuestNice

	res["user"] = gapUser * 100 / total
	res["nice"] = gapNice * 100 / total
	res["system"] = gapSystem * 100 / total
	res["idle"] = gapIdle * 100 / total
	res["iowait"] = gapIowait * 100 / total
	res["irq"] = gapIRQ * 100 / total
	res["softirq"] = gapSoftIRQ * 100 / total
	res["steal"] = gapSteal * 100 / total
	res["guest"] = gapGuest * 100 / total
	res["guest_nice"] = gapGuestNice * 100 / total

	return res, nil
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

	u := LinuxUsagePlugin{}
	plugin := mp.NewMackerelPlugin(u)
	plugin.Run()
	return 0
}
