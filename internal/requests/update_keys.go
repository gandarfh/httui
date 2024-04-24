package requests

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/terminal"
)

func (m Model) KeyActions(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Sequence(common.SetState(common.Exit_state), tea.Quit)

	case key.Matches(msg, m.keys.Detail):
		m = m.ShowRequestDetails(msg.String())
		return m, m.Detail.SetRequest(m.Requests.Current)

	case key.Matches(msg, m.keys.OpenGroup):
		m = m.OpenRequest()
		m = m.ShowRequestDetails(msg.String())
		return m, tea.Batch(
			LoadRequestsByParentId(m.parentId),
			m.Detail.SetRequest(m.Requests.Current),
		)

	case key.Matches(msg, m.keys.CloseGroup):
		m = m.BackRequest()
		return m, tea.Batch(LoadRequestsByParentId(m.previousParentId))

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
			common.OpenCommand("CREATE_WORKSPACE", "Type a name for workspace: "),
		)

	case key.Matches(msg, m.keys.Exec):
		index := m.List.Index()
		m.Requests.Current = m.Requests.List[index]
		m = m.OpenRequest()

		if m.Requests.Current.Type == "group" {
			break
		}

		return m, tea.Batch(common.SetLoading(true, "Loading..."), m.Exec())

	case key.Matches(msg, m.keys.Delete):
		index := m.List.Index()
		m.Requests.Current = m.Requests.List[index]

		return m, tea.Batch(
			common.OpenCommand("DELETE", "Type 'Y' to confirm, or any other key to cancel: "),
		)

	case key.Matches(msg, m.keys.Create):
		parentId := m.Requests.Current.ParentID

		if m.Requests.Current.Type == "group" {
			parentId = &m.Requests.Current.ID
		}

		term := terminal.NewPreview(&repositories.Request{
			ParentID: parentId,
		})

		return m, tea.Batch(term.OpenVim("Create"))

	case key.Matches(msg, m.keys.Edit):
		index := m.List.Index()
		m.Requests.Current = m.Requests.List[index]
		m = m.OpenRequest()

		if m.Requests.Current.Type == "group" {
			requests, _ := repositories.NewRequest().List(&m.Requests.Current.ID, "")

			group := struct {
				Group    repositories.Request
				Requests []repositories.Request
			}{
				Group:    m.Requests.Current,
				Requests: requests,
			}

			term := terminal.NewPreview(&group)

			return m, tea.Batch(term.OpenVim("Edit"))
		}

		term := terminal.NewPreview(&m.Requests.Current)

		return m, tea.Batch(term.OpenVim("Edit"))

	case key.Matches(msg, m.keys.Envs):
		envs, _ := repositories.NewEnvs().List(m.Workspace.ID)

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
	size := len(m.Requests.List)
	if size == 0 {
		return m
	}

	index := m.List.Index()

	switch direction {
	case "down", "j":
		index += 1
	case "up", "k":
		index -= 1
	}

	if size == index || index < 0 {
		index = m.List.Index()
	}

	m.Requests.Current = m.Requests.List[index]

	if m.Requests.Current.Type == "request" {
		m.parentId = m.Requests.Current.ParentID
	}

	return m
}

func (m Model) OpenRequest() Model {
	if len(m.Requests.List) == 0 {
		return m
	}

	index := m.List.Index()
	m.Requests.Current = m.Requests.List[index]

	if m.Requests.Current.Type == "group" {
		m.previousParentId = m.Requests.Current.ParentID
		m.parentId = &m.Requests.Current.ID
	}

	if m.Requests.Current.Type == "request" {
		m.parentId = m.Requests.Current.ParentID

		repositories.NewDefault().Update(&repositories.Default{
			RequestId: m.Requests.Current.ID,
		})
	}

	return m
}

func (m Model) BackRequest() Model {
	index := m.List.Index()
	m.Requests.Current = m.Requests.List[index]

	if m.parentId == nil {
		return m
	}

	if len(m.Requests.List) == 0 {
		m.parentId = m.previousParentId
		return m
	}

	if m.Requests.Current.ParentID != nil {
		request, _ := repositories.NewRequest().FindOne(*m.Requests.Current.ParentID)
		m.parentId = request.ParentID
		m.previousParentId = request.ParentID
	} else {
		m.previousParentId = nil
	}

	repositories.NewDefault().Update(&repositories.Default{
		RequestId: m.Requests.Current.ID,
	})

	return m
}
