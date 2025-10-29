package memory

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadMemInfo() (*MemoryInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	defer file.Close()
	memInfo := &MemoryInfo{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		valueStr := parts[1]

		switch key {
		case "MemTotal:":
			value, err := strconv.ParseUint(valueStr, 10, 64)

			if err == nil {
				memInfo.Total = value
			}

		case "MemFree:":
			value, err := strconv.ParseUint(valueStr, 10, 64)

			if err == nil {
				memInfo.Free = value
			}

		case "MemAvailable:":
			value, err := strconv.ParseUint(valueStr, 10, 64)

			if err == nil {
				memInfo.Available = value
			}

		case "Buffers:":
			value, err := strconv.ParseUint(valueStr, 10, 64)

			if err == nil {
				memInfo.Buffers = value
			}

		case "Cached:":
			value, err := strconv.ParseUint(valueStr, 10, 64)

			if err == nil {
				memInfo.Cached = value
			}

		case "SwapTotal:":
			value, err := strconv.ParseUint(valueStr, 10, 64)

			if err == nil {
				memInfo.SwapTotal = value
			}

		case "SwapFree:":
			value, err := strconv.ParseUint(valueStr, 10, 64)

			if err == nil {
				memInfo.SwapFree = value
			}

		}

	}

	return memInfo, nil
}

func ReadProcStatus() ([]Process, error) {
	var processes []Process
	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			if s, err := strconv.ParseInt(file.Name(), 10, 64); err == nil {
				processFile := fmt.Sprintf("/proc/%d/status", s)
				file, err := os.Open(processFile)
				if err != nil {
					continue
				}
				scanner := bufio.NewScanner(file)
				process := Process{}

				defer file.Close()

				for scanner.Scan() {
					line := scanner.Text()
					parts := strings.Fields(line)

					if len(parts) < 2 {
						continue
					}

					key := parts[0]
					valueStr := parts[1]

					switch key {
					case "Name:":
						process.Name = strings.Join(parts[1:], " ")

					case "Pid:":
						pid, err := strconv.ParseUint(valueStr, 10, 64)
						if err == nil {
							process.Pid = pid
						}

					case "State:":
						process.State = valueStr

					case "VmRSS:":
						vmrss, err := strconv.ParseUint(valueStr, 10, 64)
						if err == nil {
							process.VmRSS = vmrss
						}
					}

				}

				processes = append(processes, process)
			}
		}
	}

	return processes, nil

}
