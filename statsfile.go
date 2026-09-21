package main

import (
	"time"

	"github.com/monitoring-forge/saferio"
	"github.com/prometheus/procfs"
)

type stats struct {
	CPUStat procfs.CPUStat `json:"cpustat"`
	Time    int64          `json:"time"`
}

func writeStats(dir, filename string, st procfs.CPUStat) error {
	return saferio.WriteJSON(dir, filename, stats{CPUStat: st, Time: time.Now().Unix()})
}

func readStats(dir, filename string) (int64, procfs.CPUStat, error) {
	st := stats{}
	if err := saferio.ReadJSON(dir, filename, &st); err != nil {
		return 0, procfs.CPUStat{}, err
	}
	return st.Time, st.CPUStat, nil
}
