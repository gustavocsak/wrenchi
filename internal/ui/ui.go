package ui

import (
	"log"
	"time"
	"wrenchi/internal/cpu"
	"wrenchi/internal/memory"

	tea "github.com/charmbracelet/bubbletea"
)

func NewProgram() *tea.Program {
	memStats, err := memory.ReadMemInfo()
	if err != nil {
		log.Fatalf("failed to retrieve memory stats: %v", err)
	}

	cpuStats, err := cpu.ReadCPUStat()
	if err != nil {
		log.Fatalf("failed to retrieve CPU stats: %v", err)
	}

	m := model{
		memoryStats: memStats,
		cpuStats:    cpuStats,
		selected:    "memory",
		lastUpdate:  time.Now(),
	}
	return tea.NewProgram(m)
}
