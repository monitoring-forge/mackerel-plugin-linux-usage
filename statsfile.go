package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/procfs"
)

type stats struct {
	CPUStat procfs.CPUStat `json:"cpustat"`
	Time    int64          `json:"time"`
}

func fileExists(dir, filename string) bool {
	_, err := os.Stat(filepath.Join(dir, filename))
	return err == nil
}

func writeStats(dir, filename string, st procfs.CPUStat) error {
	newFile, err := os.CreateTemp(dir, "linux-usage-")
	if err != nil {
		return err
	}
	n := time.Now().Unix()
	je := json.NewEncoder(newFile)
	err = je.Encode(stats{st, n})
	if err != nil {
		newFile.Close()
		_ = os.Remove(newFile.Name())
		return err
	}

	err = newFile.Close()
	if err != nil {
		_ = os.Remove(newFile.Name())
		return err
	}

	if err := os.Rename(newFile.Name(), filepath.Join(dir, filename)); err != nil {
		_ = os.Remove(newFile.Name())
		return err
	}
	return nil
}

func readStats(dir, filename string) (int64, procfs.CPUStat, error) {
	file, err := openRD(filepath.Join(dir, filename))
	if err != nil {
		return 0, procfs.CPUStat{}, err
	}
	defer file.Close()
	st := stats{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&st)
	if err != nil {
		return 0, procfs.CPUStat{}, err
	}
	return st.Time, st.CPUStat, nil
}
