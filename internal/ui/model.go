package ui

import (
	"time"
	"wrenchi/internal/cpu"
	"wrenchi/internal/logger"
	"wrenchi/internal/memory"
	"wrenchi/internal/process"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	memoryStats   memory.MemoryInfo
	cpuStats      cpu.CPUStats
	cpuPrev       cpu.CPUStats
	cpuHistory    [][]float64
	maxHistory    int
	processTable  table.Model
	processes     []process.Process
	lastUpdate    time.Time
	width         int
	height        int
	theme         Theme
	hostname      string
	logger        *logger.Logger
	coreViewCache map[int]string
	cacheDirty    bool
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}
