package ui

import (
	"time"
	"wrenchi/internal/cpu"
	"wrenchi/internal/memory"
	"wrenchi/internal/process"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.logger.Printf("Window resized: %dx%d", msg.Width, msg.Height)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "down", "pgup", "pgdown", "home", "end":
			m.processTable, cmd = m.processTable.Update(msg)
			return m, cmd
		}

	case tickMsg:
		m.cpuPrev = m.cpuStats
		m.cacheDirty = true

		memStats, err := memory.ReadMemInfo()
		if err == nil {
			m.memoryStats = memStats
		} else {
			m.logger.Printf("Error reading memory stats: %v", err)
		}

		cpuStats, err := cpu.ReadCPUStat()
		if err == nil {
			m.cpuStats = cpuStats

			// update each core history
			for i := range cpuStats.PerCore {
				if i < len(m.cpuPrev.PerCore) && i < len(m.cpuHistory) {
					usage := cpu.CalculateCPUUsage(m.cpuPrev.PerCore[i], cpuStats.PerCore[i])

					m.cpuHistory[i] = append(m.cpuHistory[i], usage)

					if len(m.cpuHistory[i]) > m.maxHistory {
						m.cpuHistory[i] = m.cpuHistory[i][len(m.cpuHistory[i])-m.maxHistory:]
					}
				}
			}
		} else {
			m.logger.Printf("Error reading CPU stats: %v", err)
		}

		// update processes
		processes, err := process.ReadProcesses(30)
		if err == nil {
			m.processes = processes
			usage := cpu.CalculateCPUUsage(m.cpuPrev.Total, m.cpuStats.Total)
			m.processTable.SetRows(processesToRows(processes, usage))
		} else {
			m.logger.Printf("Error reading processes: %v", err)
		}

		m.lastUpdate = time.Time(msg)

		return m, tick()
	}
	return m, nil
}
