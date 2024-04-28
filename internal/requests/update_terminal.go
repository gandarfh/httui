package requests

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/terminal"
	"gorm.io/gorm"
)

func (m Model) TerminalActions(msg terminal.Finish) (Model, tea.Cmd) {
	switch msg.Category {
	case "Create":
		request := repositories.Request{}
		if err := msg.Preview.Execute(&request); err != nil {
			return m, nil
		}

		if request.Name == "" {
			return m, nil
		}

		repositories.NewRequest().Create(&request)

		m.Requests.Current = request
		m.parentId = m.Requests.Current.ParentID

		return m, tea.Batch(LoadRequestsByParentId(m.parentId))

	case "Edit":
		if m.Requests.Current.Type == "group" {
			var group = struct {
				Group    repositories.Request
				Requests []repositories.Request
			}{}

			if err := msg.Preview.Execute(&group); err != nil {
				return m, nil
			}

			for _, request := range group.Requests {
				if request.ID == 0 {
					repositories.NewRequest().Create(&request)
				}

				repositories.NewRequest().Update(&request)
			}

			repositories.NewRequest().Update(&group.Group)
			m.parentId = group.Group.ParentID

			return m, tea.Batch(LoadRequestsByParentId(m.parentId))
		}

		request := repositories.Request{}

		if err := msg.Preview.Execute(&request); err != nil {
			return m, nil
		}

		request.ID = m.Requests.Current.ID

		repositories.NewRequest().Update(&request)
		m.parentId = request.ParentID

		return m, tea.Batch(LoadRequestsByParentId(m.parentId))

	case "Envs":
		data := []map[string]any{}
		if err := msg.Preview.Execute(&data); err != nil {
			return m, nil
		}

		for _, item := range data {
			env := repositories.Env{
				WorkspaceId: m.Workspace.ID,
				Key:         item["key"].(string),
				Value:       item["value"].(string),
			}

			if item["id"] != nil {
				env.Model = gorm.Model{
					ID: uint(item["id"].(float64)),
				}
				repositories.NewEnvs().Update(&env)
			} else {
				repositories.NewEnvs().Create(&env)
			}
		}
	}

	defer msg.Preview.Close()
	if msg.Err != nil {
		return m, nil
	}

	return m, nil
}
