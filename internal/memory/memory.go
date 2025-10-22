package memory

import (
	"os"
	"bufio"
	"strings"
	"strconv"
)

func GetMemoryInfo() (*MemoryInfo, error) {
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
