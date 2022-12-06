package internal

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	}
}

func (m Model) Init() tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	cmds = append(cmds, m.spinner.Tick)

	m.pages[common.Page_Workspace].Content = workspaces.New()
	m.pages[common.Page_Resource].Content = resources.New()
	m.pages[common.Page_Env].Content = envs.New()

	for _, p := range m.pages {
		cmds = append(cmds, p.Content.Init())
	}

	m.state = loaded_state

	config, _ := m.default_repo.First()

	if config.WorkspaceId != 0 {
		cmd := common.SetPage(common.Page_Resource)
		cmds = append(cmds, cmd)

		common.CurrWorkspace, _ = m.workspace_repo.FindOne(config.WorkspaceId)
		cmd = common.ListTags(config.WorkspaceId)
		cmds = append(cmds, cmd)

	}

	if config.TagId != 0 {
		cmd := common.SetResourceTab(common.Tab_Resources)
		cmds = append(cmds, cmd)

		common.CurrTag, _ = m.tag_repo.FindOne(config.TagId)
		cmd = common.ListResources(config.TagId)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case common.Loading:
		m.loading = msg

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.state = loaded_state

		for p, i := range m.pages {
			m.pages[p].Content, cmd = i.Content.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.Page:
		common.CurrPage = msg
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
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

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	if m.state == loaded_state {
		m.pages[common.CurrPage].Content, cmd = m.pages[common.CurrPage].Content.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	content := ""

	w := m.width - (m.width / 12)
	h := m.height - (m.height / 3)

	if m.loading.Value {
		m.loading.Msg = m.spinner.View() + m.loading.Msg
	} else {
		m.loading.Msg = ""
	}

	switch m.state {
	case loaded_state:
		content = tabs.New(m.pages, common.CurrPage, w, h, m.loading)
	default:
		content = tabs.New(m.pages, common.CurrPage, w, h, m.loading)
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, styles.Container.Base.Render(content))
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
