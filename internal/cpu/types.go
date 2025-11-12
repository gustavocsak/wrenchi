package cpu

type CPUStats struct {
	CPU     string
	Total   CoreStats
	PerCore []CoreStats
	Uptime  string
}

type CoreStats struct {
	Name    string
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	IOWait  uint64
	IRQ     uint64
	SoftIRQ uint64
	MHz     float64
}

type CorePercentages struct {
	Name   string
	User   float64
	System float64
	Idle   float64
}
