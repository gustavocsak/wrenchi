package cpu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadCPUStat() (CPUStats, error) {
	statFile, err := os.Open("/proc/stat")
	if err != nil {
		return CPUStats{}, err
	}
	defer statFile.Close()
	stats, err := ParseCPUStats(statFile)
	if err != nil {
		return CPUStats{}, err
	}

	cpuInfoFile, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return stats, err
	}
	defer cpuInfoFile.Close()
	stats.CPU, err = ParseCPUInfo(cpuInfoFile)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func ParseCPUInfo(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "model name") {
			parts := strings.Split(line, ":")
			if len(parts) < 2 {
				continue
			}
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("cpu model name not found in cpuinfo")

}

func ParseCPUStats(r io.Reader) (CPUStats, error) {
	scanner := bufio.NewScanner(r)
	stats := CPUStats{}

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Fields(line)

		// cpu0 68410 36 21234 3662108 17610 4618 5058 0 0 0
		//      ^^^^^ ^^ ^^^^^ ^^^^^^^ ^^^^^ ^^^^ ^^^^
		//      user nice sys  idle   iowait irq softirq

		if len(parts) == 0 || !strings.HasPrefix(parts[0], "cpu") {
			continue
		}

		core, err := parseCoreStats(parts)
		if err != nil {
			continue
		}
		if parts[0] == "cpu" {
			stats.Total = core
		} else {
			stats.PerCore = append(stats.PerCore, core)
		}

	}

	return stats, nil
}

func parseCoreStats(parts []string) (CoreStats, error) {
	if len(parts) < 8 {
		return CoreStats{}, fmt.Errorf("not enough fields for cpu line")
	}

	core := CoreStats{Name: parts[0]}
	var err error

	core.User, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return CoreStats{}, err
	}

	core.Nice, err = strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return CoreStats{}, err
	}

	core.System, err = strconv.ParseUint(parts[3], 10, 64)
	if err != nil {
		return CoreStats{}, err
	}

	core.Idle, err = strconv.ParseUint(parts[4], 10, 64)
	if err != nil {
		return CoreStats{}, err
	}

	core.IOWait, err = strconv.ParseUint(parts[5], 10, 64)
	if err != nil {
		return CoreStats{}, err
	}

	core.IRQ, err = strconv.ParseUint(parts[6], 10, 64)
	if err != nil {
		return CoreStats{}, err
	}

	core.SoftIRQ, err = strconv.ParseUint(parts[7], 10, 64)
	if err != nil {
		return CoreStats{}, err
	}

	return core, nil
}
