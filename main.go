package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/external/database"
	"github.com/gandarfh/httui/internal"
)

func main() {
	if err := database.SqliteConnection(); err != nil {
		fmt.Println("Error to connect to database!", err)
		os.Exit(1)
	}

	m := internal.New()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
