package main

import (
	"fmt"
	"wrenchi/internal/memory"
	"wrenchi/internal/display"
)

func main() {
	fmt.Println("System Diagnostic Tool starting...")
	memInfo, error := memory.GetMemoryInfo()
	if error != nil {
		fmt.Println("error")
	}
	display.PrintMemoryInfo(memInfo)
}
