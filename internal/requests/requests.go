package requests

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/terminal"
	"gorm.io/gorm"
)

type Model struct {
	width            int
	height           int
	state            common.State
	help             help.Model
	keys             KeyMap
	request_list     list.Model
	request_detail   ModelDetail
	title            string
	filter           string
	parentId         *uint
	previousParentId *uint
}

var (
	divider = lipgloss.NewStyle().MarginLeft(1).Border(lipgloss.NormalBorder(), false, true, false, false)
)

func New() common.Component {

	m := Model{
		state:          common.Start_state,
		request_list:   NewRequestList(),
		request_detail: NewDetail(),
		help:           help.New(),
		keys:           keys,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case Result:
		if msg.Err != nil {
			term := terminal.NewPreview(&msg.Err)
			return m, tea.Batch(common.SetLoading(false), term.OpenVim("Exec"))
		}

		term := terminal.NewPreview(&msg.Response)
		return m, tea.Batch(common.SetLoading(false), term.OpenVim("Exec"))

	case common.State:
		m.state = msg
		switch msg {
		case common.Start_state:
		case common.Loaded_state:
			m.request_detail.Request = common.CurrRequest
			m.request_detail.Preview = fmt.Sprintf("%s - %s", common.CurrRequest.Method, common.CurrRequest.Endpoint)

			m.parentId = common.CurrRequest.ParentID
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.OpenGroup):
			m = m.OpenRequest()

		case key.Matches(msg, m.keys.CloseGroup):
			m = m.BackRequest()

		case key.Matches(msg, m.keys.Filter):
			return m, tea.Batch(
				common.OpenCommand("FILTER", ""),
			)

		case key.Matches(msg, m.keys.Workspace):
			return m, tea.Batch(
				common.OpenCommand("SET_WORKSPACE", ""),
			)

		case key.Matches(msg, m.keys.Exec):
			index := m.request_list.Index()
			common.CurrRequest = common.ListOfRequests[index]
			m = m.OpenRequest()

			if common.CurrRequest.Type == "group" {
				break
			}

			return m, tea.Batch(common.SetLoading(true, "Loading..."), m.Exec())

		case key.Matches(msg, m.keys.Delete):
			index := m.request_list.Index()
			common.CurrRequest = common.ListOfRequests[index]

			return m, tea.Batch(
				common.OpenCommand("DELETE", "Type 'Y' to confirm, or any other key to cancel: "),
			)

		case key.Matches(msg, m.keys.Create):
			term := terminal.NewPreview(&repositories.Request{
				ParentID: common.CurrRequest.ParentID,
			})

			return m, tea.Batch(term.OpenVim("Create"))

		case key.Matches(msg, m.keys.Edit):
			index := m.request_list.Index()
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

	case common.Tab:
		common.CurrTab = msg

	case common.List:
		common.ListOfRequests = msg.Requests

		m.request_list.SetItems(m.RequestOfList())
		m.request_list, cmd = m.request_list.Update(msg)
		cmds = append(cmds, cmd)

	case terminal.Finish:
		switch msg.Category {
		case "Create":
			request := repositories.Request{}
			msg.Preview.Execute(&request)
			repositories.NewRequest().Create(&request)

			common.CurrRequest = request
			m.parentId = common.CurrRequest.ParentID

			return m, tea.Batch(common.ListRequests(*common.CurrRequest.ParentID))

		case "Edit":
			msg.Preview.Execute(&common.CurrRequest)
			repositories.NewRequest().Update(&common.CurrRequest)
			m.parentId = common.CurrRequest.ParentID

			return m, tea.Batch(common.ListRequests(*common.CurrRequest.ParentID))

		case "Envs":
			data := []map[string]any{}
			msg.Preview.Execute(&data)

			for _, env := range data {
				env := repositories.Env{
					Model: gorm.Model{
						ID: uint(env["id"].(float64)),
					},
					WorkspaceId: common.CurrWorkspace.ID,
					Key:         env["key"].(string),
					Value:       env["value"].(string),
				}

				repositories.NewEnvs().Update(&env)
			}
		}

		defer msg.Preview.Close()
		if msg.Err != nil {
			return m, nil
		}

	case common.CommandClose:
		switch msg.Category {
		case "DELETE":
			if strings.ToUpper(msg.Value) == "Y" {
				repositories.NewRequest().Delete(common.CurrRequest.ID)
			}

		case "FILTER":
			m.filter = msg.Value

		case "SET_WORKSPACE":
			workspace := repositories.Workspace{}
			repositories.NewWorkspace().Sql.Model(&workspace).Where("name LIKE ?", "%"+msg.Value+"%").First(&workspace)
			common.CurrWorkspace = workspace

			repositories.NewDefault().Update(&repositories.Default{
				WorkspaceId: workspace.ID,
			})

			cmd = common.SetEnvironment(workspace.Name)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m = m.WindowSize(msg)
	}

	m.request_detail, cmd = m.request_detail.Update(msg)
	cmds = append(cmds, cmd)

	m.request_list.SetItems(m.RequestOfList())
	m.request_list, cmd = m.request_list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Right,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.request_list.View(),
			divider.Height(m.height).String(),
			m.request_detail.View(),
		),
	)
}
