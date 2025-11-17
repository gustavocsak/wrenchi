package process

type Process struct {
	PID     int
	Name    string
	State   string
	CPUTime uint64
	RSS     uint64
	Command string
}

type ProcessInfo struct {
	PID        int
	Name       string
	CPUPercent float64
	MemPercent float64
	Command    string
}
