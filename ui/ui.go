package ui

import (
	"github/gandarfh/httui/ui/components/header"
	"github/gandarfh/httui/ui/home"
	"github/gandarfh/httui/ui/requests"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle     = lipgloss.NewStyle().Padding(1, 3).Background(lipgloss.Color("#090909"))
	contentStyle = lipgloss.NewStyle().Align(lipgloss.Left).Background(lipgloss.Color("#090909"))
)

type page int

const (
	homePage page = iota
	requestPage
)

func (p *page) Get() string {
	return []string{"homeView", "requestView"}[*p]
}

func (p *page) Set(new page) {
	*p = new
}

type Model struct {
	home     home.Model
	requests requests.Model
	header   header.Model
	currPage page
	Width    int
	Height   int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.home.Init(), tea.EnterAltScreen)
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		homeCmd     tea.Cmd
		requestsCmd tea.Cmd
		cmds        []tea.Cmd
	)

	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)

	case tea.KeyMsg:
		switch msg.String() {

		case "H":
			m.currPage.Set(homePage)
			return m, nil

		case "L":
			m.currPage.Set(requestPage)
			return m, nil

		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	m.home, homeCmd = m.home.Update(message)
	m.requests, requestsCmd = m.requests.Update(message)

	cmds = append(cmds, homeCmd, requestsCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}

	tabs := home.Tabs(m.Width-9, int(m.currPage))

	s.WriteString(tabs + "\n\n")

	if m.currPage == homePage {
		s.WriteString(m.home.View())
	}

	if m.currPage == requestPage {
		s.WriteString(m.requests.View())
	}

	return appStyle.
		Width(m.Width).
		Height(m.Height).
		Render(m.header.View() + "\n\n" + contentStyle.Width(m.Width-7).Render(s.String()))
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	m.Width = msg.Width
	m.Height = msg.Height
}
