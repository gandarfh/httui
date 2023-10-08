package requests

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/common"
)

func (m Model) CommandsActions(msg common.CommandClose) (Model, tea.Cmd) {
	switch msg.Category {
	case "DELETE":
		if strings.ToUpper(msg.Value) == "Y" {
			repositories.NewRequest().Delete(common.CurrRequest.ID)
		}

	case "FILTER":
		m.filter = msg.Value
		return m, tea.Batch(common.ListRequests(nil))

	case "CREATE_WORKSPACE":
		workspace := repositories.Workspace{Name: msg.Value}
		repositories.NewWorkspace().Create(&workspace)
		common.CurrWorkspace = workspace

		repositories.NewDefault().Update(&repositories.Default{
			WorkspaceId: workspace.ID,
		})

		return m, common.SetEnvironment(workspace.Name)

	case "SET_WORKSPACE":
		workspace := repositories.Workspace{}
		repositories.NewWorkspace().Sql.Model(&workspace).Where("name LIKE ?", "%"+msg.Value+"%").First(&workspace)
		common.CurrWorkspace = workspace

		repositories.NewDefault().Update(&repositories.Default{
			WorkspaceId: workspace.ID,
		})

		return m, common.SetEnvironment(workspace.Name)
	}

	return m, nil
}
