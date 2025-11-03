package cpu

import (
	"io"
	"os"
)

func ReadMemInfo() (CPU, error) {
	file, err := os.Open("/proc/cpuinfo")

	if err != nil {
		return CPU{}, err
	}
	defer file.Close()
	return ParseCpuInfo(file)
}

func ParseCpuInfo(r io.Reader) (CPU, error) {
	//TODO: read /proc/cpuinfo and populate the CPU struct
	return CPU{}, nil
}
