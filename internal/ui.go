package internal

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/maid-san/internal/command"
	"github.com/gandarfh/maid-san/internal/envs"
	"github.com/gandarfh/maid-san/internal/repositories"
	"github.com/gandarfh/maid-san/internal/resources"
	"github.com/gandarfh/maid-san/internal/workspaces"
	"github.com/gandarfh/maid-san/pkg/common"
	"github.com/gandarfh/maid-san/pkg/styles"
	"github.com/gandarfh/maid-san/pkg/tabs"
)

type state int

const (
	start_state state = iota
	error_state
	loaded_state
)

type Model struct {
	default_repo   *repositories.DefaultsRepo
	workspace_repo *repositories.WorkspacesRepo
	tag_repo       *repositories.TagsRepo
	resource_repo  *repositories.ResourcesRepo
	pages          tabs.Contents
	spinner        spinner.Model
	loading        common.Loading
	width          int
	height         int
	state          state
	command_page   common.Component
	command_active bool
}

func New() Model {
	var (
		default_repo, _   = repositories.NewDefault()
		workspace_repo, _ = repositories.NewWorkspace()
		tag_repo, _       = repositories.NewTag()
		resource_repo, _  = repositories.NewResource()
	)

	pages := tabs.Contents{
		{Tab: "Workspaces"},
		{Tab: "Resouces"},
		{Tab: "Environments"},
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().MarginLeft(2).Foreground(styles.DefaultTheme.PrimaryText)

	return Model{
		default_repo:   default_repo,
		workspace_repo: workspace_repo,
		tag_repo:       tag_repo,
		resource_repo:  resource_repo,
		pages:          pages,
		state:          start_state,
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

	m.pages[common.Page_Workspace].Content = workspaces.New()
	m.pages[common.Page_Resource].Content = resources.New()
	m.pages[common.Page_Env].Content = envs.New()

	for _, p := range m.pages {
		cmds = append(cmds, p.Content.Init())
	}

	config, _ := m.default_repo.First()

	if config.WorkspaceId != 0 {
		cmd = common.SetWorkspace(config.WorkspaceId)
		cmds = append(cmds, cmd)

		cmd = common.ListTags(config.WorkspaceId)
		cmds = append(cmds, cmd)

		cmd = common.SetPage(common.Page_Resource)
		cmds = append(cmds, cmd)
	}

	if config.TagId != 0 {
		cmd = common.SetTag(config.TagId)
		cmds = append(cmds, cmd)

		common.CurrTag, _ = m.tag_repo.FindOne(config.TagId)
		cmd = common.ListResources(config.TagId)
		cmds = append(cmds, cmd)

		cmd = common.SetTab(common.Tab_Resources)
		cmds = append(cmds, cmd)
	}

	m.state = loaded_state

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case common.CommandClose:
		m.command_active = false

	case common.Command:
		m.command_active = msg.Active

	case common.Loading:
		m.loading = msg

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 5

		m.state = loaded_state

		for p, i := range m.pages {
			m.pages[p].Content, cmd = i.Content.Update(msg)
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
			case "right", "l":
				common.CurrPage = min(common.CurrPage+1, len(m.pages)-1)
				return m, nil
			case "left", "h":
				common.CurrPage = max(common.CurrPage-1, 0)
				return m, nil
			}
		}
	}

	if m.state == loaded_state {
		if m.command_active {
			m.command_page, cmd = m.command_page.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.pages[common.CurrPage].Content, cmd = m.pages[common.CurrPage].Content.Update(msg)
			cmds = append(cmds, cmd)

		}

		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
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
			lipgloss.Top, tabs.New(m.pages, common.CurrPage, w, m.height, m.loading),
			m.command_page.View(),
		)
	} else {
		content = tabs.New(m.pages, common.CurrPage, w, m.height, m.loading)
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		styles.Container.Base.Render(content),
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
