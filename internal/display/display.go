package display

import (
	"fmt"
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

	memUsedRatio := float64(info.Total - info.Available) / float64(info.Total)
	fmt.Printf("%-20s %.2f%%\n", "Memory Usage:", memUsedRatio * 100)

	fmt.Println("\n---Swap Usage---")
	fmt.Printf(format, "Total Swap:", formatMemory(info.SwapTotal))
	fmt.Printf(format, "Free Swap:", formatMemory(info.SwapFree))
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
