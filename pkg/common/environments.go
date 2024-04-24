package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
)

type Environment struct {
	Workspace repositories.Workspace
}

func SetWorkspace(workspaceId uint) tea.Cmd {
	return func() tea.Msg {
		workspace_repo := repositories.NewWorkspace()
		workspace, _ := workspace_repo.FindOne(workspaceId)

		return Environment{
			Workspace: workspace,
		}
	}
}
