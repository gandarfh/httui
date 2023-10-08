package requests

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/terminal"
)

func (m Model) KeyActions(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Detail):
		m = m.ShowRequestDetails(msg.String())

	case key.Matches(msg, m.keys.OpenGroup):
		m = m.OpenRequest()
		return m, tea.Batch(common.ListRequests(m.parentId))

	case key.Matches(msg, m.keys.CloseGroup):
		m = m.BackRequest()
		return m, tea.Batch(common.ListRequests(m.previousParentId))

	case key.Matches(msg, m.keys.Filter):
		return m, tea.Batch(
			common.OpenCommand("FILTER", ""),
		)

	case key.Matches(msg, m.keys.SetWorkspace):
		return m, tea.Batch(
			common.OpenCommand("SET_WORKSPACE", ""),
		)

	case key.Matches(msg, m.keys.CreateWorkspace):
		return m, tea.Batch(
			common.OpenCommand("CREATE_WORKSPACE", ""),
		)

	case key.Matches(msg, m.keys.Exec):
		index := m.List.Index()
		common.CurrRequest = common.ListOfRequests[index]
		m = m.OpenRequest()

		if common.CurrRequest.Type == "group" {
			break
		}

		return m, tea.Batch(common.SetLoading(true, "Loading..."), m.Exec())

	case key.Matches(msg, m.keys.Delete):
		index := m.List.Index()
		common.CurrRequest = common.ListOfRequests[index]

		return m, tea.Batch(
			common.OpenCommand("DELETE", "Type 'Y' to confirm, or any other key to cancel: "),
		)

	case key.Matches(msg, m.keys.Create):
		parentId := common.CurrRequest.ParentID

		if common.CurrRequest.Type == "group" {
			parentId = &common.CurrRequest.ID
		}

		term := terminal.NewPreview(&repositories.Request{
			ParentID: parentId,
		})

		return m, tea.Batch(term.OpenVim("Create"))

	case key.Matches(msg, m.keys.Edit):
		index := m.List.Index()
		common.CurrRequest = common.ListOfRequests[index]
		m = m.OpenRequest()

		term := terminal.NewPreview(&common.CurrRequest)

		return m, tea.Batch(term.OpenVim("Edit"))

	case key.Matches(msg, m.keys.Envs):
		envs, _ := repositories.NewEnvs().List(common.CurrWorkspace.ID)

		data := []map[string]any{}
		for _, env := range envs {
			data = append(data, map[string]any{"id": env.ID, "key": env.Key, "value": env.Value})
		}

		term := terminal.NewPreview(&data)

		return m, tea.Batch(term.OpenVim("Envs"))
	}

	return m, nil
}

func (m Model) ShowRequestDetails(direction string) Model {
	size := len(common.ListOfRequests)
	if size == 0 {
		return m
	}

	index := m.List.Cursor()

	switch direction {
	case "down", "j":
		index += 1
	case "up", "k":
		index -= 1
	}

	if size == index || index < 0 {
		index = m.List.Cursor()
	}

	common.CurrRequest = common.ListOfRequests[index]

	if common.CurrRequest.Type == "request" {
		m.parentId = common.CurrRequest.ParentID

		m.detail.Request = common.CurrRequest
		m.detail.Preview = fmt.Sprintf("%s - %s", common.CurrRequest.Method, common.CurrRequest.Endpoint)
	}

	return m
}

func (m Model) OpenRequest() Model {
	if len(common.ListOfRequests) == 0 {
		return m
	}

	index := m.List.Index()
	common.CurrRequest = common.ListOfRequests[index]

	if common.CurrRequest.Type == "group" {
		m.parentId = &common.CurrRequest.ID
		m.previousParentId = common.CurrRequest.ParentID
	}

	if common.CurrRequest.Type == "request" {
		m.parentId = common.CurrRequest.ParentID

		m.detail.Request = common.CurrRequest
		m.detail.Preview = fmt.Sprintf("%s - %s", common.CurrRequest.Method, common.CurrRequest.Endpoint)

		repositories.NewDefault().Update(&repositories.Default{
			RequestId: common.CurrRequest.ID,
		})
	}

	return m
}

func (m Model) BackRequest() Model {
	if m.parentId == nil {
		return m
	}

	if len(common.ListOfRequests) == 0 {
		m.parentId = m.previousParentId
		return m
	}

	group, _ := repositories.NewRequest().FindOne(*m.parentId)
	common.CurrRequest = *group

	m.parentId = common.CurrRequest.ParentID

	repositories.NewDefault().Update(&repositories.Default{
		RequestId: common.CurrRequest.ID,
	})

	return m
}
