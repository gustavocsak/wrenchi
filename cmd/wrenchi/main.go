package main

import "fmt"
import "wrenchi/internal/memory"

func main() {
	fmt.Println("System Diagnostic Tool starting...")
	mockData, error := memory.GetMemoryInfo()
	fmt.Println("Success:", mockData)
	fmt.Println("error:", error)
}
