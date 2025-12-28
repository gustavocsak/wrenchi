package ui

import (
	"log"
	"time"
	"wrenchi/internal/cpu"
	"wrenchi/internal/logger"
	"wrenchi/internal/memory"
	"wrenchi/internal/process"
	"wrenchi/internal/system"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	// initial process list
	processes, err := process.ReadProcesses(30)
	if err != nil {
		log.Printf("warning: failed to read processes: %v", err)
		processes = []process.Process{}
	}

	columns := []table.Column{
		{Title: "PID", Width: 8},
		{Title: "Name", Width: 20},
		{Title: "CPU%", Width: 8},
		{Title: "MEM%", Width: 8},
		{Title: "Command", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(processesToRows(processes, 0)),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	theme := DefaultTheme()
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.Border).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(theme.Accent).
		Bold(true)
	t.SetStyles(s)

	cpuHistory := make([][]float64, len(cpuStats.PerCore))
	for i := range cpuHistory {
		cpuHistory[i] = make([]float64, 0, 30)
	}

	l, err := logger.New("wrenchi.log")
	if err != nil {
		log.Fatalf("unable to initialize logger")
	}

	m := model{
		memoryStats:   memStats,
		cpuStats:      cpuStats,
		cpuPrev:       cpuStats,
		cpuHistory:    cpuHistory,
		maxHistory:    30,
		processTable:  t,
		processes:     processes,
		lastUpdate:    time.Now(),
		theme:         theme,
		hostname:      system.Hostname(),
		logger:        l,
		coreViewCache: make(map[int]string),
		cacheDirty:    true,
	}

	return tea.NewProgram(m, tea.WithAltScreen())
}
