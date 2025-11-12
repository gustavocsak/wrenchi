package main

import (
	"log"
	"wrenchi/internal/ui"
)

func main() {
	p := ui.NewProgram()
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
