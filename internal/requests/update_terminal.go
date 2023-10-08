package requests

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/terminal"
	"gorm.io/gorm"
)

func (m Model) TerminalActions(msg terminal.Finish) (Model, tea.Cmd) {
	switch msg.Category {
	case "Create":
		request := repositories.Request{}
		msg.Preview.Execute(&request)
		repositories.NewRequest().Create(&request)

		common.CurrRequest = request
		m.parentId = common.CurrRequest.ParentID

		return m, tea.Batch(common.ListRequests(common.CurrRequest.ParentID))

	case "Edit":
		if common.CurrRequest.Type == "group" {
			var group = struct {
				Group    repositories.Request
				Requests []repositories.Request
			}{}

			msg.Preview.Execute(&group)

			for _, request := range group.Requests {
				repositories.NewRequest().Update(&request)
			}

			repositories.NewRequest().Update(&group.Group)
			m.parentId = group.Group.ParentID

			return m, tea.Batch(common.ListRequests(common.CurrRequest.ParentID))
		}

		request := repositories.Request{}
		msg.Preview.Execute(&request)
		request.ID = common.CurrRequest.ID

		repositories.NewRequest().Update(&request)
		m.parentId = common.CurrRequest.ParentID

		return m, tea.Batch(common.ListRequests(common.CurrRequest.ParentID))

	case "Envs":
		data := []map[string]any{}
		msg.Preview.Execute(&data)

		for _, item := range data {
			env := repositories.Env{
				WorkspaceId: common.CurrWorkspace.ID,
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
