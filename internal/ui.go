package internal

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/command"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/internal/requests"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/styles"
	"github.com/gandarfh/httui/pkg/tabs"
)

type Model struct {
	width          int
	height         int
	workspaceName  string
	pages          tabs.Contents
	spinner        spinner.Model
	state          common.State
	command_active bool
	environment    common.Environment
	loading        common.Loading
	command_page   common.Component
	request_repo   *repositories.RequestsRepo
}

func New() Model {
	pages := tabs.Contents{
		{Tab: "Requests", Content: requests.New()},
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().MarginLeft(2).Foreground(styles.DefaultTheme.PrimaryText)

	return Model{
		environment:    common.Environment{Name: common.CurrWorkspace.Name},
		request_repo:   repositories.NewRequest(),
		pages:          pages,
		state:          common.Start_state,
		spinner:        s,
		command_page:   command.New(),
		command_active: false,
	}
}

func (m Model) Init() tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	cmds = append(cmds, m.spinner.Tick)
	cmds = append(cmds, m.command_page.Init())

	for _, p := range m.pages {
		cmds = append(cmds, p.Content.Init())
	}

	cmd = common.SetState(common.Start_state)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case common.State:
		m.state = msg
		switch msg {
		case common.Error_state:
		case common.Loaded_state:
		case common.Start_state:
			conf, _ := repositories.NewDefault().First()

			request, _ := repositories.NewRequest().FindOne(conf.RequestId)
			common.CurrRequest = *request

			workspace, _ := repositories.NewWorkspace().FindOne(conf.WorkspaceId)
			common.CurrWorkspace = workspace
			m.environment.Name = workspace.Name

			cmd = common.SetState(common.Loaded_state)
			cmds = append(cmds, cmd)
		default:
		}

	case common.CommandClose:
		m.command_active = false

	case common.Command:
		m.command_active = msg.Active

	case common.Environment:
		m.environment = msg

	case common.Loading:
		m.loading = msg

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 5

		for p, i := range m.pages {
			content, cmd := i.Content.Update(msg)
			m.pages[p].Content = content.(common.Component)

			cmds = append(cmds, cmd)
		}

	case common.Page:
		common.CurrPage = msg
		return m, nil

	case tea.KeyMsg:
		if !m.command_active {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	if m.state == common.Loaded_state {
		if m.command_active {
			content, cmd := m.command_page.Update(msg)
			m.command_page = content.(common.Component)

			cmds = append(cmds, cmd)
		} else {
			content, cmd := m.pages[common.CurrPage].Content.Update(msg)
			m.pages[common.CurrPage].Content = content.(common.Component)

			cmds = append(cmds, cmd)

		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	content := ""

	if m.loading.Value {
		m.loading.Msg = m.spinner.View() + m.loading.Msg
	} else {
		m.loading.Msg = ""
	}

	w := m.width - 2

	if m.command_active {
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			tabs.New(m.pages, int(common.CurrPage), w, m.height, m.loading, m.environment),
			m.command_page.View(),
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Right,
			tabs.New(m.pages, int(common.CurrPage), w, m.height, m.loading, m.environment),
			lipgloss.NewStyle().MarginRight(2).Render(m.pages[common.CurrPage].Content.Help()),
		)
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		styles.Container.Base.Render(content),
	)
}
