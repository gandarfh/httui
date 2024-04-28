package requests

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/terminal"
)

type UpdateRequestDefault repositories.Request

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case RequestsData:
		m.Requests = msg

		cmds = append(cmds, m.List.SetItems(m.RequestOfList()))

		if len(m.Requests.List) > 0 {
			index := m.List.Index()

			if len(m.Requests.List) <= index {
				index = 0
			}

			m.Requests.Current = m.Requests.List[index]
			cmds = append(cmds, m.Detail.SetRequest(m.Requests.Current))
			cmds = append(cmds, func() tea.Msg { return UpdateRequestDefault(m.Requests.Current) })
		}

	case repositories.Workspace:
		m.Workspace = msg
		m.List.Title = fmt.Sprintf("[%s]", msg.Name)

	case repositories.Default:
		m.Configs = msg

	case UpdateRequestDefault:
		if m.Requests.Current.ID == msg.ID {
			repositories.NewDefault().Update(&repositories.Default{
				RequestId: m.Requests.Current.ID,
			})
		}

	case Result:
		if msg.Err != nil {
			term := terminal.NewPreview(&msg.Err)
			return m, tea.Batch(common.SetLoading(false), term.OpenVim("Exec"))
		}

		term := terminal.NewPreview(&msg.Response)
		return m, tea.Batch(common.SetLoading(false), term.OpenVim("Exec"))

	case common.State:
		m, cmd = m.StateActions(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		if !m.command_active {
			m, cmd = m.KeyActions(msg)
			cmds = append(cmds, cmd)
		}

	case common.Environment:
		m.Workspace = msg.Workspace
		m.List.Title = fmt.Sprintf("[%s]", msg.Workspace.Name)

		cmd = m.Detail.SetWorkspace(repositories.Workspace(m.Workspace))
		cmds = append(cmds, cmd)

	case spinner.TickMsg:
		if m.loading.Value {
			m.spinner, cmd = m.spinner.Update(msg)
			m.List.Title = fmt.Sprintf("[%s]", m.Workspace.Name) + m.spinner.View()
			cmds = append(cmds, cmd)
		}

	case common.Loading:
		m.loading = msg

		if m.loading.Value {
			cmds = append(cmds, m.spinner.Tick)
		}

	case terminal.Finish:
		m.List.Title = fmt.Sprintf("[%s]", m.Workspace.Name)
		m, cmd = m.TerminalActions(msg)

		cmds = append(cmds, cmd)

	case common.Command:
		m.command_active = msg.Active

	case common.CommandClose:
		m.command_active = false
		m, cmd = m.CommandsActions(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width

		m.List.SetHeight(m.Height / 2)
		m.List.SetWidth(m.Width / 5)

		m.Detail.Width = m.Width - m.List.Width() - 2
		m.Detail.Height = m.Height - 9
	}

	if m.command_active {
		content, cmd := m.command_bar.Update(msg)
		m.command_bar = content.(common.Component)

		return m, cmd
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	m.Detail, cmd = m.Detail.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) RequestOfList() []list.Item {
	list := []list.Item{}
	w := m.List.Width() - 2

	for _, i := range m.Requests.List {
		list = append(list, RequestItem{i.Name, i.Type, w})
	}

	return list
}
