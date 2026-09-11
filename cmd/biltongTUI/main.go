package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	tui "github.com/Tillter2998/biltongTUI"
)

func main() {
	p := tea.NewProgram(tui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
