package memory

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"strconv"
	"strings"
)

func ReadMemInfo() (*MemoryInfo, error) {
	file, err := os.Open("/proc/meminfo")

	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseMemInfo(file)
}

func ParseMemInfo(r io.Reader) (*MemoryInfo, error) {
	memInfo := &MemoryInfo{}
	scanner := bufio.NewScanner(r)

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
	dirs, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	var processes []Process

	for _, entry := range dirs {
		if entry.IsDir() {
			if s, err := strconv.ParseUint(entry.Name(), 10, 64); err == nil {
				p, err := openProcessFile(s)
				if err != nil {
					continue
				}
				processes = append(processes, p)
			}
		}
	}

	return processes, nil
}

func openProcessFile(pid uint64) (Process, error) {
	processFileName := fmt.Sprintf("/proc/%d/status", pid)

	file, err := os.Open(processFileName)
	if err != nil {
		return Process{}, err
	}
	defer file.Close()

	p := ParseProcess(file)
	return p, nil
}

func ParseProcess(r io.Reader) Process {
	scanner := bufio.NewScanner(r)
	process := Process{}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "Name:":
			process.Name = strings.Join(parts[1:], " ")

		case "Pid:":
			pid, err := strconv.ParseUint(value, 10, 64)
			if err == nil {
				process.Pid = pid
			}

		case "State:":
			process.State = value

		case "VmRSS:":
			vmrss, err := strconv.ParseUint(value, 10, 64)
			if err == nil {
				process.VmRSS = vmrss
			}

		case "Uid:":
			id := strings.TrimSpace(parts[1])
			user, err := user.LookupId(id)
			if err == nil {
				process.User = user.Username
			}
			uid, err := ParseId(parts[1:])
			if err == nil {
				process.Uid = uid
			}

		case "Gid:":
			group, err := user.LookupGroupId(parts[1])
			if err == nil {
				process.Group = group.Name
			}

			gid, err := ParseId(parts[1:])
			if err == nil {
				process.Gid = gid
			}
		}

	}

	return process

}

func ParseId(ids []string) (IDs, error) {
	var rid IDs
	if len(ids) < 4 || len(ids) > 5 {
		return IDs{}, fmt.Errorf("Unable to parse ids from process file")
	}

	for i, v := range ids {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return IDs{}, err
		}
		switch i {
		case 0:
			rid.Effective = uint(id)
		case 1:
			rid.Real = uint(id)
		case 2:
			rid.Saved = uint(id)
		case 3:
			rid.FS = uint(id)
		}
	}

	return rid, nil
}
