package requests

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/pkg/terminal"
	"gorm.io/gorm"
)

func (m Model) TerminalActions(msg terminal.Finish) (Model, tea.Cmd) {
	switch msg.Category {
	case "Create":
		request := offline.Request{}
		if err := msg.Preview.Execute(&request); err != nil {
			return m, nil
		}

		if request.Name == "" {
			return m, nil
		}

		offline.NewRequest().Create(&request)

		m.Requests.Current = request
		m.parentId = m.Requests.Current.ParentID

		return m, tea.Batch(LoadRequestsByParentId(m.parentId))

	case "Edit":
		if m.Requests.Current.Type == "group" {
			var group = struct {
				Group    offline.Request
				Requests []offline.Request
			}{}

			if err := msg.Preview.Execute(&group); err != nil {
				return m, nil
			}

			for _, request := range group.Requests {
				if request.ID == 0 {
					offline.NewRequest().Create(&request)
				}

				offline.NewRequest().Update(&request)
			}

			offline.NewRequest().Update(&group.Group)
			m.parentId = group.Group.ParentID

			return m, tea.Batch(LoadRequestsByParentId(m.parentId))
		}

		request := offline.Request{}

		if err := msg.Preview.Execute(&request); err != nil {
			return m, nil
		}

		request.ID = m.Requests.Current.ID

		offline.NewRequest().Update(&request)
		m.parentId = request.ParentID

		return m, tea.Batch(LoadRequestsByParentId(m.parentId))

	case "Envs":
		data := []map[string]any{}
		if err := msg.Preview.Execute(&data); err != nil {
			return m, nil
		}

		for _, item := range data {
			env := offline.Env{
				WorkspaceId: m.Workspace.ID,
				Key:         item["key"].(string),
				Value:       item["value"].(string),
			}

			if item["id"] != nil {
				env.Model = gorm.Model{
					ID: uint(item["id"].(float64)),
				}
				offline.NewEnvs().Update(&env)
			} else {
				offline.NewEnvs().Create(&env)
			}
		}
	}

	defer msg.Preview.Close()
	if msg.Err != nil {
		return m, nil
	}

	return m, nil
}
