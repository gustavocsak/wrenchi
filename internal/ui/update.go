package ui

import (
	"time"
	"wrenchi/internal/cpu"
	"wrenchi/internal/memory"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tickMsg:
		// Refresh memory stats
		memStats, err := memory.ReadMemInfo()
		if err == nil {
			m.memoryStats = memStats
		}

		// Refresh CPU stats
		cpuStats, err := cpu.ReadCPUStat()
		if err == nil {
			m.cpuStats = cpuStats
		}

		// Update last refresh time
		m.lastUpdate = time.Time(msg)

		// Schedule next tick
		return m, tick()
	}
	return m, nil
}
