package requests

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/command"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/internal/requests/details"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/styles"
)

func LoadDefault() tea.Msg {
	config, _ := repositories.NewDefault().First()
	return *config
}

func LoadWorspace() tea.Msg {
	config, _ := repositories.NewDefault().First()
	workspace, _ := repositories.NewWorkspace().FindOne(config.WorkspaceId)
	return workspace
}

type RequestsData struct {
	List    []repositories.Request
	Current repositories.Request
}

func LoadRequests() tea.Msg {
	config, _ := repositories.NewDefault().First()
	request, _ := repositories.NewRequest().FindOne(config.RequestId)
	requests, _ := repositories.NewRequest().List(request.ParentID, "")

	return RequestsData{
		Current: *request,
		List:    requests,
	}
}

func LoadRequestsByParentId(parentId *uint) tea.Cmd {
	return func() tea.Msg {
		requests, _ := repositories.NewRequest().List(parentId, "")
		return RequestsData{
			List: requests,
		}
	}
}

func LoadRequestsByFilter(filter string) tea.Cmd {
	return func() tea.Msg {
		requests, _ := repositories.NewRequest().List(nil, filter)
		return RequestsData{
			List: requests,
		}
	}
}

type Model struct {
	Detail           details.Model
	title            string
	filter           string
	parentId         *uint
	previousParentId *uint
	command_active   bool
	keys             KeyMap
	help             help.Model
	List             list.Model
	spinner          spinner.Model
	command_bar      common.Component
	loading          common.Loading
	state            common.State
	Width            int
	Height           int
	Requests         RequestsData
	Configs          repositories.Default
	Workspace        repositories.Workspace
}

var (
	divider = lipgloss.NewStyle().MarginLeft(1).Border(lipgloss.NormalBorder(), false, true, false, false)
)

func New() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().MarginLeft(2).Foreground(styles.DefaultTheme.PrimaryText)

	m := Model{
		Width:          0,
		Height:         0,
		state:          common.Start_state,
		List:           NewRequestList(),
		Detail:         details.New(),
		help:           help.New(),
		keys:           keys,
		spinner:        s,
		command_bar:    command.New(),
		command_active: false,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Sequence(
		LoadDefault,
		LoadWorspace,
		LoadRequests,
		m.command_bar.Init(),
		m.Detail.Init(),
		common.SetState(common.Start_state),
	)
}

var (
	Base = lipgloss.NewStyle()
)

func (m Model) View() string {
	requestsPage := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Height(m.Height-3).Width(m.List.Width()-4).Render(m.List.View()),
		m.Detail.View(),
	)

	if m.state != common.Loaded_state {
		return ""
	}

	footer := ""

	if m.command_active {
		footer = m.command_bar.View()
	} else {
		footer = lipgloss.NewStyle().Render(m.Help())
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		requestsPage,
		footer,
	)

	return Base.
		Height(m.Height).
		Width(m.Width).
		Render(content)
}
