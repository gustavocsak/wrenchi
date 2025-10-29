package display

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"sort"
	"strconv"
	"wrenchi/internal/memory"
)

func PrintMemoryInfo(info *memory.MemoryInfo) {
	const format = "%-20s %s\n"
	fmt.Println("\n---Memory Usage---")
	fmt.Printf(format, "Total Memory:", formatMemory(info.Total))
	fmt.Printf(format, "Available Memory:", formatMemory(info.Available))
	fmt.Printf(format, "Free Memory:", formatMemory(info.Free))
	fmt.Printf(format, "Buffers:", formatMemory(info.Buffers))
	fmt.Printf(format, "Cached:", formatMemory(info.Cached))

	memUsedRatio := (float64(info.Total-info.Available) / float64(info.Total)) * 100.0
	memBar := progress.New(progress.WithWidth(40))
	memBar.FullColor = "41"
	memTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("91"))
	fmt.Printf(format, "Memory Usage:", memTextStyle.Render(memBar.ViewAs(memUsedRatio/100.0)))

	fmt.Println("\n---Swap Usage---")
	fmt.Printf(format, "Total Swap:", formatMemory(info.SwapTotal))
	fmt.Printf(format, "Free Swap:", formatMemory(info.SwapFree))

	swapUsedKB := float64(info.SwapTotal - info.SwapFree)
	var swapUsedRatio float64
	if info.SwapTotal > 0 {
		swapUsedRatio = (swapUsedKB / float64(info.SwapTotal)) * 100
	}

	swapBar := progress.New(progress.WithWidth(40))
	swapBar.FullColor = "31"
	swapTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("91"))
	swapBar.SetPercent(swapUsedRatio / 100.0)
	fmt.Printf(format, "Swap Usage:", swapTextStyle.Render(swapBar.ViewAs(swapUsedRatio/100.0)))
	fmt.Println()
}

func PrintProcesses(processes []memory.Process, limit int) {
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].VmRSS > processes[j].VmRSS
	})
	const format = "%-6d %-30s %-3s %-10s\n"
	var baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	columns := []table.Column{
		{Title: "Pid", Width: 6},
		{Title: "Name", Width: 20},
		{Title: "State", Width: 6},
		{Title: "Memory Usage", Width: 15},
	}

	rows := []table.Row{}

	for _, p := range processes[:limit] {
		pid := strconv.FormatUint(p.Pid, 10)
		vmrssStr := formatMemory(p.VmRSS)
		rows = append(rows, table.Row{pid, p.Name, p.State, vmrssStr})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(limit+2),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	t.SetStyles(s)

	fmt.Println(baseStyle.Render(t.View()))
}

func formatMemory(kb uint64) string {
	const (
		mb = 1024
		gb = 1024 * 1024
	)

	if kb > gb {
		return fmt.Sprintf("%.2f GB", float64(kb)/float64(gb))
	}

	if kb > mb {
		return fmt.Sprintf("%.2f MB", float64(kb)/float64(mb))
	}

	return fmt.Sprintf("%d KB", kb)
}
