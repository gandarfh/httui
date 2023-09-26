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
	Width            int
	Height           int
	state            common.State
	help             help.Model
	keys             KeyMap
	List             list.Model
	request_detail   ModelDetail
	title            string
	filter           string
	parentId         *uint
	previousParentId *uint
}

var (
	divider = lipgloss.NewStyle().MarginLeft(1).Border(lipgloss.NormalBorder(), false, true, false, false)
)

func New() Model {
	m := Model{
		Width:          100,
		state:          common.Start_state,
		List:           NewRequestList(),
		request_detail: NewDetail(),
		help:           help.New(),
		keys:           keys,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
			m.List.Title = fmt.Sprintf("[%s]", common.CurrWorkspace.Name)

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

	case common.Tab:
		common.CurrTab = msg

	case common.Environment:
		m.List.Title = fmt.Sprintf("[%s]", msg.Name)

	case common.List:
		common.ListOfRequests = msg.Requests

		m.List.SetItems(m.RequestOfList())
		m.List, cmd = m.List.Update(msg)
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

	case common.CommandClose:
		switch msg.Category {
		case "DELETE":
			if strings.ToUpper(msg.Value) == "Y" {
				repositories.NewRequest().Delete(common.CurrRequest.ID)
			}

		case "FILTER":
			m.filter = msg.Value

		case "CREATE_WORKSPACE":
			workspace := repositories.Workspace{Name: msg.Value}
			repositories.NewWorkspace().Create(&workspace)
			common.CurrWorkspace = workspace

			repositories.NewDefault().Update(&repositories.Default{
				WorkspaceId: workspace.ID,
			})

			cmd = common.SetEnvironment(workspace.Name)
			cmds = append(cmds, cmd)

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

	m.List.SetItems(m.RequestOfList())
	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.List.View(),
		divider.Height(m.Height-1).String(),
		m.request_detail.View(),
	)
}
