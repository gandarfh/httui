package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/config"
	"github.com/gandarfh/httui/internal/login"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/internal/requests"
)

func init() {
	if err := offline.SqliteConnection(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if config.Config.Settings.Logging {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
}

func App() {
	m := requests.New()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func Login() {
	m := login.New()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
