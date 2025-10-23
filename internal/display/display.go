package display

import (
	"fmt"
	"wrenchi/internal/memory"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/progress"
)

func PrintMemoryInfo(info *memory.MemoryInfo) {
	const format = "%-20s %s\n"
	fmt.Println("\n---Memory Usage---")
	fmt.Printf(format, "Total Memory:", formatMemory(info.Total))
	fmt.Printf(format, "Available Memory:", formatMemory(info.Available))
	fmt.Printf(format, "Free Memory:", formatMemory(info.Free))
	fmt.Printf(format, "Buffers:", formatMemory(info.Buffers))
	fmt.Printf(format, "Cached:", formatMemory(info.Cached))

	memUsedRatio := (float64(info.Total - info.Available) / float64(info.Total)) * 100.0
	memBar := progress.New(progress.WithWidth(40))
	memBar.FullColor = "41"
	memTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("91"))
	fmt.Printf(format, "Memory Usage:", memTextStyle.Render(memBar.ViewAs(memUsedRatio / 100.0)))



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
	fmt.Printf(format, "Swap Usage:", swapTextStyle.Render(swapBar.ViewAs(swapUsedRatio / 100.0)))
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
