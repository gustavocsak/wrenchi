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

	memUsed := info.Total - info.Available
	fmt.Printf(format, "Used Memory:", formatMemory(memUsed))
	fmt.Printf(format, "Free Memory:", formatMemory(info.Free))
	fmt.Printf(format, "Buffers:", formatMemory(info.Buffers))
	fmt.Printf(format, "Cached:", formatMemory(info.Cached))
	memUsedRatio := (float64(memUsed) / float64(info.Total)) * 100.0
	memBar := progress.New(progress.WithWidth(32))
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

	swapBar := progress.New(progress.WithWidth(32))
	swapBar.FullColor = "31"
	swapTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("91"))
	swapBar.SetPercent(swapUsedRatio / 100.0)
	fmt.Printf(format, "Swap Usage:", swapTextStyle.Render(swapBar.ViewAs(swapUsedRatio/100.0)))
	fmt.Println()
}

func PrintProcesses(processes []memory.Process, limit int, info *memory.MemoryInfo) {
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].VmRSS > processes[j].VmRSS
	})
	var baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	columns := []table.Column{
		{Title: "Pid", Width: 7},
		{Title: "Name", Width: 20},
		{Title: "user", Width: 10},
		{Title: "State", Width: 6},
		{Title: "Memory Usage", Width: 15},
		{Title: "%MEM", Width: 6},
	}

	rows := []table.Row{}

	for _, p := range processes[:limit] {
		pid := strconv.FormatUint(p.Pid, 10)
		vmrssStr := formatMemory(p.VmRSS)
		memPercentage := (float64(p.VmRSS) / float64(info.Total)) * 100.0
		memPercentageStr := fmt.Sprintf("%.2f%%", memPercentage)
		rows = append(rows, table.Row{pid, p.Name, p.User, p.State, vmrssStr, memPercentageStr})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(limit+2),
	)

	s := table.DefaultStyles()
	s.Selected = lipgloss.NewStyle()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Padding(0, 1)

	t.SetStyles(s)
	fmt.Println(baseStyle.Render(t.View()))
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
