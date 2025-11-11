package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"wrenchi/internal/cpu"
	"wrenchi/internal/display"
	"wrenchi/internal/memory"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must provide subcommand\n")
		os.Exit(1)
	}

	subCommand := os.Args[1]

	switch subCommand {
	case "memory":
		memoryCmd := flag.NewFlagSet("memory", flag.ExitOnError)
		showTopMem := memoryCmd.Int("t", 0, "Show top t processes by memory use")
		memoryCmd.Parse(os.Args[2:])
		memInfo, err := memory.ReadMemInfo()
		if err != nil {
			fmt.Println(err)
		}
		display.PrintMemoryInfo(memInfo)
		if *showTopMem != 0 {
			processes, err := memory.ReadProcStatus()
			if err != nil {
				fmt.Println(err)
			}
			display.PrintProcesses(processes, *showTopMem, memInfo)
		}
	case "cpu":
		stats, err := cpu.ReadCPUStat()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(stats)
	default:
		log.Fatalf("Invalid subcommand: %s", subCommand)
		os.Exit(1)
	}

}
