package main

import (
	"github/gandarfh/httui/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	p := tea.NewProgram(
		ui.Model{},
		tea.WithAltScreen(),
	)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
