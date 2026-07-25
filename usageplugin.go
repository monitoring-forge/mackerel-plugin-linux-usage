package main

import (
	"fmt"
	"maps"
	"os"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/prometheus/procfs"
)

var tooOldDuration = 600.0 // seconds

type LinuxUsagePlugin struct {
	workDir string
}

func (u LinuxUsagePlugin) tempfilePath() string {
	uid := os.Geteuid()
	return fmt.Sprintf("mackerel-plugin-linux-usage-%d", uid)
}

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

func (u LinuxUsagePlugin) gaugeMetrics(pf procfs.FS) (map[string]float64, error) {
	res := map[string]float64{}
	st, err := pf.Stat()
	if err != nil {
		return res, err
	}
	loadavg, err := pf.LoadAvg()
	if err != nil {
		return res, err
	}
	procs, err := pf.AllProcs()
	if err != nil {
		return res, err
	}
	selffs, err := pf.Self()
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

	cores := float64(len(st.CPU))
	if cores == 0 {
		cores = 1
	}

	res["loadavg1"] = loadavg.Load1 / cores
	res["loadavg5"] = loadavg.Load5 / cores
	res["loadavg15"] = loadavg.Load15 / cores
	res["all"] = totalProcs
	res["running"] = procRunning
	if psnmp.Tcp.ActiveOpens != nil {
		res["active"] = *psnmp.Tcp.ActiveOpens
	}
	if psnmp.Tcp.PassiveOpens != nil {
		res["passive"] = *psnmp.Tcp.PassiveOpens
	}
	if pnetstat.ListenOverflows != nil {
		res["overflows"] = *pnetstat.ListenOverflows
	}
	if pnetstat.ListenDrops != nil {
		res["drops"] = *pnetstat.ListenDrops
	}

	return res, nil
}

func minZero(a float64) float64 {
	if a < 0 {
		return 0
	}
	return a
}

func (u LinuxUsagePlugin) cpuMetrics(pf procfs.FS) (map[string]float64, error) {
	res := map[string]float64{}

	cur, err := pf.Stat()
	if err != nil {
		return res, err
	}

	path := u.tempfilePath()
	curCPU := cur.CPUTotal

	defer func() {
		if writeErr := writeStats(u.workDir, path, curCPU); writeErr != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to save stats to %s: %v\n", path, writeErr)
		}
	}()

	if !fileExists(u.workDir, path) {
		fmt.Fprintf(os.Stderr, "Notice: first time execution command\n")
		return res, nil
	}

	prevTime, prevCPU, err := readStats(u.workDir, path)
	if err != nil {
		return res, err
	}
	if prevTime == 0 {
		return res, fmt.Errorf("failed to get previous time")
	}
	now := time.Now().Unix()
	timeDiff := float64(now - prevTime)
	if timeDiff > tooOldDuration {
		return res, fmt.Errorf("too long duration")
	}

	var total float64
	// User
	gapUser := minZero(float64(curCPU.User) - float64(prevCPU.User))
	total += gapUser
	// Nice
	gapNice := minZero(float64(curCPU.Nice) - float64(prevCPU.Nice))
	total += gapNice
	// System
	gapSystem := minZero(float64(curCPU.System) - float64(prevCPU.System))
	total += gapSystem
	// Idle
	gapIdle := minZero(float64(curCPU.Idle) - float64(prevCPU.Idle))
	total += gapIdle
	// Iowait
	gapIowait := minZero(float64(curCPU.Iowait) - float64(prevCPU.Iowait))
	total += gapIowait
	// Irq
	gapIRQ := minZero(float64(curCPU.IRQ) - float64(prevCPU.IRQ))
	total += gapIRQ
	// SoftIRQ
	gapSoftIRQ := minZero(float64(curCPU.SoftIRQ) - float64(prevCPU.SoftIRQ))
	total += gapSoftIRQ
	// Steal
	gapSteal := minZero(float64(curCPU.Steal) - float64(prevCPU.Steal))
	total += gapSteal
	// Guest
	gapGuest := minZero(float64(curCPU.Guest) - float64(prevCPU.Guest))
	total += gapGuest
	// GuestNice
	gapGuestNice := minZero(float64(curCPU.GuestNice) - float64(prevCPU.GuestNice))
	total += gapGuestNice

	if total == 0 {
		fmt.Fprintf(os.Stderr, "Notice: System CPU counter seems to be unchanged\n")
		return res, nil
	}

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

func (u LinuxUsagePlugin) FetchMetrics() (map[string]float64, error) {
	pf, err := procfs.NewDefaultFS()
	if err != nil {
		return nil, err
	}
	res, err := u.gaugeMetrics(pf)
	if err != nil {
		return nil, err
	}

	cpu, err := u.cpuMetrics(pf)
	if err != nil {
		return nil, err
	}

	maps.Copy(res, cpu)
	return res, nil

}
