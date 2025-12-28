package process

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// reads all processes from /proc and returns them sorted by CPU usage
func ReadProcesses(limit int) ([]Process, error) {
	totalMemKB, err := getTotalMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get total memory: %w", err)
	}

	// page size for RSS calculation
	pageSize := uint64(os.Getpagesize())

	procDirs, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		return nil, fmt.Errorf("failed to glob /proc: %w", err)
	}

	var processes []Process
	for _, dir := range procDirs {
		proc, err := readProcess(dir)
		if err != nil {
			continue
		}
		processes = append(processes, proc)
	}

	// later: implement calculations using delta
	var totalCPUTime uint64
	for _, p := range processes {
		totalCPUTime += p.CPUTime
	}

	for i := range len(processes) {
		p := &processes[i]
		var cpuPercent float64
		if totalCPUTime > 0 {
			cpuPercent = (float64(p.CPUTime) / float64(totalCPUTime)) * 100.0
		}

		// RSS is the number of resident pages being used
		// RSS * pageSize will give the total bytes being used
		memKB := (p.RSS * pageSize) / 1024
		memPercent := (float64(memKB) / float64(totalMemKB)) * 100.0

		p.CPUPercent = cpuPercent
		p.MemPercent = memPercent
		p.MemKB = memKB
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPUPercent > processes[j].CPUPercent
	})
	if limit > 0 && limit < len(processes) {
		processes = processes[:limit]
	}

	return processes, nil
}

// reads a single process from /proc/[pid]
func readProcess(procDir string) (Process, error) {
	// dir name = pid
	pidStr := filepath.Base(procDir)
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return Process{}, fmt.Errorf("invalid PID: %w", err)
	}

	// read /proc/[pid]/stat
	statPath := filepath.Join(procDir, "stat")
	statData, err := os.ReadFile(statPath)
	if err != nil {
		return Process{}, err
	}

	statStr := string(statData)

	commStart := strings.Index(statStr, "(")
	commEnd := strings.LastIndex(statStr, ")")
	if commStart == -1 || commEnd == -1 {
		return Process{}, fmt.Errorf("invalid stat format")
	}

	name := statStr[commStart+1 : commEnd]

	fields := strings.Fields(statStr[commEnd+2:])
	if len(fields) < 21 {
		return Process{}, fmt.Errorf("insufficient fields in stat")
	}

	// parse fields from stat
	state := fields[0]

	utime, _ := strconv.ParseUint(fields[11], 10, 64)
	stime, _ := strconv.ParseUint(fields[12], 10, 64)
	cpuTime := utime + stime

	rss, _ := strconv.ParseUint(fields[21], 10, 64)

	cmdlinePath := filepath.Join(procDir, "cmdline")
	cmdlineData, err := os.ReadFile(cmdlinePath)
	command := name
	if err == nil && len(cmdlineData) > 0 {
		command = strings.ReplaceAll(string(cmdlineData), "\x00", " ")
		command = strings.TrimSpace(command)
		if command == "" {
			command = fmt.Sprintf("[%s]", name)
		}
	}

	return Process{
		PID:     pid,
		Name:    name,
		State:   state,
		CPUTime: cpuTime,
		RSS:     rss,
		Command: command,
	}, nil
}

// reads total memory from /proc/meminfo
func getTotalMemory() (uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				kb, err := strconv.ParseUint(fields[1], 10, 64)
				if err == nil {
					return kb, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("MemTotal not found in /proc/meminfo")
}
