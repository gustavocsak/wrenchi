package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"wrenchi/internal/cpu"
	"wrenchi/internal/memory"
)

func (m model) View() string {
	var s strings.Builder

	// Header with title and last update time
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(0, 1)

	updateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	s.WriteString(titleStyle.Render("WRENCHI SYSTEM MONITOR"))
	s.WriteString("  ")
	if !m.lastUpdate.IsZero() {
		s.WriteString(updateStyle.Render(
			fmt.Sprintf("Last update: %s", m.lastUpdate.Format("15:04:05"))))
	}
	s.WriteString("\n")

	s.WriteString(displayMemory(m.memoryStats))
	s.WriteString("\n")
	s.WriteString(displayCPU(m.cpuStats))

	// Footer with help text
	s.WriteString("\n")
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	s.WriteString(helpStyle.Render("Press 'q' to quit • Updates every 2 seconds"))

	return s.String()
}

func displayMemory(mem memory.MemoryInfo) string {
	const format = "%-20s %s\n"
	var s strings.Builder

	// Memory section header
	s.WriteString("\n=== MEMORY USAGE ===\n\n")

	// Memory stats
	s.WriteString(fmt.Sprintf(format, "Total Memory:", formatMemory(mem.Total)))
	s.WriteString(fmt.Sprintf(format, "Available Memory:", formatMemory(mem.Available)))

	memUsed := mem.Total - mem.Available
	s.WriteString(fmt.Sprintf(format, "Used Memory:", formatMemory(memUsed)))
	s.WriteString(fmt.Sprintf(format, "Free Memory:", formatMemory(mem.Free)))
	s.WriteString(fmt.Sprintf(format, "Buffers:", formatMemory(mem.Buffers)))
	s.WriteString(fmt.Sprintf(format, "Cached:", formatMemory(mem.Cached)))

	// Memory usage progress bar
	memUsedRatio := (float64(memUsed) / float64(mem.Total)) * 100.0
	memBar := progress.New(progress.WithWidth(32))
	memBar.FullColor = "41"
	memTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("91"))
	s.WriteString(fmt.Sprintf(format, "Memory Usage:",
		memTextStyle.Render(memBar.ViewAs(memUsedRatio/100.0))))

	// Swap section header
	s.WriteString("\n=== SWAP USAGE ===\n\n")
	s.WriteString(fmt.Sprintf(format, "Total Swap:", formatMemory(mem.SwapTotal)))
	s.WriteString(fmt.Sprintf(format, "Free Swap:", formatMemory(mem.SwapFree)))

	// Swap usage calculation and progress bar
	swapUsedKB := float64(mem.SwapTotal - mem.SwapFree)
	var swapUsedRatio float64
	if mem.SwapTotal > 0 {
		swapUsedRatio = (swapUsedKB / float64(mem.SwapTotal)) * 100
	}

	swapBar := progress.New(progress.WithWidth(32))
	swapBar.FullColor = "31"
	swapTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("91"))
	swapBar.SetPercent(swapUsedRatio / 100.0)
	s.WriteString(fmt.Sprintf(format, "Swap Usage:",
		swapTextStyle.Render(swapBar.ViewAs(swapUsedRatio/100.0))))

	s.WriteString("\n")

	return s.String()
}

func displayCPU(cpu cpu.CPUStats) string {
	var s strings.Builder

	// CPU section header
	s.WriteString("\n=== CPU INFORMATION ===\n\n")

	// CPU model name
	if cpu.CPU != "" {
		cpuNameStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)
		s.WriteString(fmt.Sprintf("Model: %s\n\n", cpuNameStyle.Render(cpu.CPU)))
	}

	// Number of cores
	coreCountStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	s.WriteString(fmt.Sprintf("Cores: %s\n\n",
		coreCountStyle.Render(fmt.Sprintf("%d", len(cpu.PerCore)))))

	// Per-core information header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)

	s.WriteString(headerStyle.Render(fmt.Sprintf("%-8s  %-12s  %-10s\n",
		"CORE", "FREQUENCY", "STATUS")))
	s.WriteString(headerStyle.Render(strings.Repeat("─", 40)))
	s.WriteString("\n")

	// Display each core
	for _, core := range cpu.PerCore {
		coreName := core.Name
		if coreName == "" {
			coreName = "unknown"
		}

		// Format MHz
		var mhzStr string
		if core.MHz > 0 {
			mhzStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
			mhzStr = mhzStyle.Render(fmt.Sprintf("%.2f MHz", core.MHz))
		} else {
			mhzStr = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("N/A")
		}

		// Status indicator (based on idle time - simple heuristic for now)
		var status string
		total := core.User + core.Nice + core.System + core.Idle + core.IOWait + core.IRQ + core.SoftIRQ
		if total > 0 {
			idlePercent := (float64(core.Idle) / float64(total)) * 100
			if idlePercent > 90 {
				status = lipgloss.NewStyle().
					Foreground(lipgloss.Color("34")).
					Render("● Idle")
			} else if idlePercent > 50 {
				status = lipgloss.NewStyle().
					Foreground(lipgloss.Color("226")).
					Render("● Active")
			} else {
				status = lipgloss.NewStyle().
					Foreground(lipgloss.Color("196")).
					Render("● Busy")
			}
		} else {
			status = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("● Unknown")
		}

		s.WriteString(fmt.Sprintf("%-8s  %-12s  %s\n", coreName, mhzStr, status))
	}

	return s.String()
}

func formatMemory(kb uint64) string {
	const (
		mb = 1024
		gb = 1024 * 1024
	)

	const format = "%-8.2f %-3s"

	if kb > gb {
		return fmt.Sprintf(format, float64(kb)/float64(gb), "GB")
	}

	if kb > mb {
		return fmt.Sprintf(format, float64(kb)/float64(mb), "MB")
	}

	return fmt.Sprintf(format, float64(kb), "KB")
}
