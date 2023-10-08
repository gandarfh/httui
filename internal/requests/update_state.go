package requests

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/common"
)

func (m Model) StateActions(msg common.State) (Model, tea.Cmd) {
	m.state = msg

	switch msg {
	case common.Start_state:
		conf, _ := repositories.NewDefault().First()

		request, _ := repositories.NewRequest().FindOne(conf.RequestId)
		common.CurrRequest = *request

		m.parentId = common.CurrRequest.ParentID

		workspace, _ := repositories.NewWorkspace().FindOne(conf.WorkspaceId)
		common.CurrWorkspace = workspace

		return m, common.SetState(common.Loaded_state)

	case common.Loaded_state:
		m.List.SetItems(m.RequestOfList())

		m.detail.Request = common.CurrRequest
		m.detail.Preview = fmt.Sprintf("%s - %s", common.CurrRequest.Method, common.CurrRequest.Endpoint)
		m.List.Title = fmt.Sprintf("[%s]", common.CurrWorkspace.Name)

		m.parentId = common.CurrRequest.ParentID
	}

	return m, nil
}
