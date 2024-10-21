package requests

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/pkg/common"
)

func (m Model) CommandsActions(msg common.CommandClose) (Model, tea.Cmd) {
	switch msg.Category {
	case "DELETE":
		if strings.ToUpper(msg.Value) == "Y" {
			offline.NewRequest().Delete(m.Requests.Current.ID)
		}

		return m, tea.Batch(LoadRequestsByParentId(m.Requests.Current.ParentID))

	case "FILTER":
		m.filter = msg.Value
		m.List.CursorTop()
		return m, tea.Batch(LoadRequestsByFilter(m.filter))

	case "CREATE_WORKSPACE":
		if msg.Value == "" {
			return m, nil
		}

		workspace := offline.Workspace{Name: msg.Value}
		offline.NewWorkspace().Create(&workspace)
		m.Workspace = workspace

		offline.NewDefault().Update(offline.Default{
			WorkspaceId: workspace.ID,
		})

		return m, common.SetWorkspace(workspace.ID)

	case "SET_WORKSPACE":
		if msg.Value == "" {
			return m, nil
		}

		workspace := offline.Workspace{}
		offline.NewWorkspace().Sql.Model(&workspace).Where("name LIKE ?", "%"+msg.Value+"%").First(&workspace)
		m.Workspace = workspace

		offline.NewDefault().Update(offline.Default{
			WorkspaceId: workspace.ID,
		})

		return m, common.SetWorkspace(workspace.ID)
	}

	return m, nil
}
