package main

import (
	"fmt"
	"wrenchi/internal/display"
	"wrenchi/internal/memory"
)

func main() {
	fmt.Println("System Diagnostic Tool starting...")
	memInfo, error := memory.ReadMemInfo()
	if error != nil {
		fmt.Println("error")
	}
	display.PrintMemoryInfo(memInfo)

	processes, err := memory.ReadProcStatus()
	if err != nil {
		fmt.Println("error")
	}
	// probably would pass flags here?
	display.PrintProcesses(processes)
}
