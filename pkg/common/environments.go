package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories/offline"
)

type Environment struct {
	Workspace offline.Workspace
}

func SetWorkspace(workspaceId uint) tea.Cmd {
	return func() tea.Msg {
		workspace_repo := offline.NewWorkspace()
		workspace, _ := workspace_repo.FindOne(workspaceId)

		return Environment{
			Workspace: workspace,
		}
	}
}
