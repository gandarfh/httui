package requests

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/terminal"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	log.Printf("[%T]: %v\n", msg, msg)
	switch msg := msg.(type) {
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

	case common.Tab:
		common.CurrTab = msg

	case common.Environment:
		m.List.Title = fmt.Sprintf("[%s]", msg.Name)

	case common.List:
		common.ListOfRequests = msg.Requests

		m.List.SetItems(m.RequestOfList())
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)

	case spinner.TickMsg:
		if m.loading.Value {
			m.spinner, cmd = m.spinner.Update(msg)
			m.List.Title = fmt.Sprintf("[%s]", common.CurrWorkspace.Name) + m.spinner.View()

			return m, cmd
		}

	case common.Loading:
		m.loading = msg

		if m.loading.Value {
			cmds = append(cmds, m.spinner.Tick)
		}

	case terminal.Finish:
		m.List.Title = fmt.Sprintf("[%s]", common.CurrWorkspace.Name)
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
		m.Width = msg.Width + 1

		m.List.SetHeight(m.Height/2 - 2)
		m.List.SetWidth(m.Width / 5)
		m.detail.Height = ((m.Height) - 9)
		m.detail.Width = m.Width - m.List.Width() + 1
	}

	if m.command_active {
		content, cmd := m.command_bar.Update(msg)
		m.command_bar = content.(common.Component)

		return m, cmd
	}

	m.detail, cmd = m.detail.Update(msg)
	cmds = append(cmds, cmd)

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
