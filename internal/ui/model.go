package ui

import (
	"time"
	"wrenchi/internal/cpu"
	"wrenchi/internal/memory"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	memoryStats memory.MemoryInfo
	cpuStats    cpu.CPUStats
	selected    string
	lastUpdate  time.Time
}

// tickMsg is sent on every timer tick
type tickMsg time.Time

// tick returns a Cmd that waits for the refresh interval and then sends a tickMsg
func tick() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	// Start the ticker when the program starts
	return tick()
}
