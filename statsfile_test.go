//go:build !windows

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/monitoring-forge/saferio"
	"github.com/prometheus/procfs"
)

func TestWriteStats(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "test-stats.json"

	cpuStat := procfs.CPUStat{
		User:      100.0,
		Nice:      10.0,
		System:    50.0,
		Idle:      1000.0,
		Iowait:    5.0,
		IRQ:       0.0,
		SoftIRQ:   0.0,
		Steal:     0.0,
		Guest:     0.0,
		GuestNice: 0.0,
	}

	err := writeStats(tmpDir, filename, cpuStat)
	if err != nil {
		t.Fatalf("writeStats failed: %v", err)
	}

	// Verify file exists
	if !saferio.FileExists(tmpDir, filename) {
		t.Error("writeStats should create the target file")
	}

	// Verify file content
	storedTime, storedCPUStat, err := readStats(tmpDir, filename)
	if err != nil {
		t.Fatalf("readStats failed: %v", err)
	}

	now := time.Now().Unix()
	if storedTime < now-5 || storedTime > now+5 {
		t.Errorf("stored time %d is not close to current time %d", storedTime, now)
	}

	if storedCPUStat.User != cpuStat.User {
		t.Errorf("expected User %f, got %f", cpuStat.User, storedCPUStat.User)
	}
	if storedCPUStat.Nice != cpuStat.Nice {
		t.Errorf("expected Nice %f, got %f", cpuStat.Nice, storedCPUStat.Nice)
	}
	if storedCPUStat.System != cpuStat.System {
		t.Errorf("expected System %f, got %f", cpuStat.System, storedCPUStat.System)
	}
	if storedCPUStat.Idle != cpuStat.Idle {
		t.Errorf("expected Idle %f, got %f", cpuStat.Idle, storedCPUStat.Idle)
	}
	if storedCPUStat.Iowait != cpuStat.Iowait {
		t.Errorf("expected Iowait %f, got %f", cpuStat.Iowait, storedCPUStat.Iowait)
	}
	if storedCPUStat.Steal != cpuStat.Steal {
		t.Errorf("expected Steal %f, got %f", cpuStat.Steal, storedCPUStat.Steal)
	}
	if storedCPUStat.Guest != cpuStat.Guest {
		t.Errorf("expected Guest %f, got %f", cpuStat.Guest, storedCPUStat.Guest)
	}
	if storedCPUStat.GuestNice != cpuStat.GuestNice {
		t.Errorf("expected GuestNice %f, got %f", cpuStat.GuestNice, storedCPUStat.GuestNice)
	}
}

func TestWriteStatsInvalidDir(t *testing.T) {
	err := writeStats("/nonexistent/directory", "test.json", procfs.CPUStat{})
	if err == nil {
		t.Error("writeStats should return an error for invalid directory")
	}
}

func TestReadStatsNonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	_, _, err := readStats(tmpDir, "nonexistent.json")
	if err == nil {
		t.Error("readStats should return an error for non-existent file")
	}
}

func TestWriteReadStatsMultipleTimes(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "test-stats-multiple.json"

	// Write first stats
	cpuStat1 := procfs.CPUStat{
		User: 100.0,
		Idle: 1000.0,
	}
	err := writeStats(tmpDir, filename, cpuStat1)
	if err != nil {
		t.Fatalf("first writeStats failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	// Write second stats
	cpuStat2 := procfs.CPUStat{
		User: 200.0,
		Idle: 2000.0,
	}
	err = writeStats(tmpDir, filename, cpuStat2)
	if err != nil {
		t.Fatalf("second writeStats failed: %v", err)
	}

	// Read and verify second stats
	_, storedCPUStat, err := readStats(tmpDir, filename)
	if err != nil {
		t.Fatalf("readStats failed: %v", err)
	}

	if storedCPUStat.User != cpuStat2.User {
		t.Errorf("expected User %f, got %f", cpuStat2.User, storedCPUStat.User)
	}
	if storedCPUStat.Idle != cpuStat2.Idle {
		t.Errorf("expected Idle %f, got %f", cpuStat2.Idle, storedCPUStat.Idle)
	}
}

func TestStatsJSONMarshalUnmarshal(t *testing.T) {
	st := stats{
		CPUStat: procfs.CPUStat{
			User:      100.0,
			Nice:      10.0,
			System:    50.0,
			Idle:      1000.0,
			Iowait:    5.0,
			IRQ:       0.0,
			SoftIRQ:   0.0,
			Steal:     0.0,
			Guest:     0.0,
			GuestNice: 0.0,
		},
		Time: 1234567890,
	}

	// Marshal
	data, err := json.Marshal(st)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	// Unmarshal
	var decoded stats
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Time != st.Time {
		t.Errorf("expected time %d, got %d", st.Time, decoded.Time)
	}
	if decoded.CPUStat.User != st.CPUStat.User {
		t.Errorf("expected User %f, got %f", st.CPUStat.User, decoded.CPUStat.User)
	}
}

func TestCorruptedStatsFile(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "corrupted.json"

	// Write corrupted content
	corruptFile := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(corruptFile, []byte("not valid json {{{"), 0644); err != nil {
		t.Fatalf("failed to create corrupted file: %v", err)
	}

	_, _, err := readStats(tmpDir, filename)
	if err == nil {
		t.Error("readStats should return an error for corrupted file content")
	}
}

func TestReadStatsLegacyFormat(t *testing.T) {
	dir := t.TempDir()
	// State written before the saferio migration must remain readable.
	data := []byte(`{"cpustat":{"User":100,"Idle":900},"time":1234567890}`)
	if err := os.WriteFile(filepath.Join(dir, "legacy.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	timestamp, cpu, err := readStats(dir, "legacy.json")
	if err != nil {
		t.Fatal(err)
	}
	if timestamp != 1234567890 || cpu.User != 100 || cpu.Idle != 900 {
		t.Fatalf("unexpected legacy state: time=%d cpu=%+v", timestamp, cpu)
	}
}

func TestStatsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.json")
	original := []byte(`{"cpustat":{"User":100},"time":1234567890}`)
	if err := os.WriteFile(target, original, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(dir, "state.json")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readStats(dir, "state.json"); err == nil {
		t.Fatal("readStats must reject symlinks")
	}
	if err := writeStats(dir, "state.json", procfs.CPUStat{User: 200}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(original) {
		t.Fatal("writeStats modified the symlink target")
	}
	if _, cpu, err := readStats(dir, "state.json"); err != nil || cpu.User != 200 {
		t.Fatalf("replacement state: cpu=%+v err=%v", cpu, err)
	}
}
