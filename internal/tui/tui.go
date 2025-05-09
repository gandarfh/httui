package tui

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/config"
	"github.com/gandarfh/httui/internal/login"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/internal/repositories/sync"
	"github.com/gandarfh/httui/internal/requests"
	"github.com/gandarfh/httui/internal/services"
)

func init() {
	if err := offline.SqliteConnection(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

var (
	program           tea.Program
	mqttCmd           = sync.MQTTConnect(&program)
	syncRequestsCmd   = sync.SyncRequests(&program)
	syncWorkspacesCmd = sync.SyncWorkspaces(&program)
	syncResponsesCmd  = sync.SyncResponses(&program)
)

func App() {
	configDir, _ := os.UserHomeDir()
	path := filepath.Join(configDir, config.AppDir, "debug.log")

	f, err := tea.LogToFile(path, "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	services.DatasourceStart()

	m := requests.New(
		mqttCmd,
		syncRequestsCmd,
		syncWorkspacesCmd,
		syncResponsesCmd,
	)

	program = *tea.NewProgram(m, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func Login() {
	configDir, _ := os.UserHomeDir()
	path := filepath.Join(configDir, config.AppDir, "debug.log")

	f, err := tea.LogToFile(path, "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	services.DatasourceStart()

	m := login.New()
	program = *tea.NewProgram(m)

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
