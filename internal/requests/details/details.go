package details

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/utils"
)

var (
	styleConfig   = glamour.DarkStyleConfig
	title_style   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	preview_style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
)

type Model struct {
	Width     int
	Height    int
	Request   repositories.Request
	Workspace repositories.Workspace
}

func New() Model {
	return Model{}
}

func (m Model) SetWorkspace(w repositories.Workspace) tea.Cmd {
	return func() tea.Msg {
		return w
	}
}

func (m Model) SetRequest(r repositories.Request) tea.Cmd {
	return func() tea.Msg {
		return r
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case repositories.Request:
		m.Request = msg

	case repositories.Workspace:
		m.Workspace = msg
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	params_box := lipgloss.
		NewStyle().
		Height(m.Height).
		Width(m.Width / 3).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			m.Params(),
			lipgloss.NewStyle().Width(m.Width/3).Border(lipgloss.NormalBorder(), true, false, false, false).String(),
			m.Headers(),
		))

	content_row := lipgloss.NewStyle().Width(m.Width).Height(m.Height).Border(lipgloss.RoundedBorder()).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.Body(),
		lipgloss.NewStyle().Height(m.Height).Width(1).Border(lipgloss.NormalBorder(), false, true, false, false).String(),
		params_box,
	))

	container := lipgloss.NewStyle().Padding(0, 1).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.Title(),
		m.Preview(),
		content_row,
	))

	return container

}

func replaceStringsInJSON(input string, replaceFunc func(string) string) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return "", err
	}

	data = modifyData(data, replaceFunc)

	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func modifyData(value interface{}, replaceFunc func(string) string) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		if len(v) > 7 {
			v = limitMapSize(v, 7)
		}
		return processMap(v, replaceFunc)
	case []interface{}:
		if len(v) > 3 {
			v = v[:3]
		}
		return processSlice(v, replaceFunc)
	case string:
		return replaceFunc(v)
	default:
		return v
	}
}

func limitMapSize(m map[string]interface{}, max int) map[string]interface{} {
	newMap := make(map[string]interface{})
	count := 0
	for key, value := range m {
		if count >= max {
			break
		}
		newMap[key] = value
		count++
	}
	return newMap
}

func processMap(m map[string]interface{}, replaceFunc func(string) string) map[string]interface{} {
	for key, val := range m {
		m[key] = modifyData(val, replaceFunc)
	}
	return m
}

func processSlice(s []interface{}, replaceFunc func(string) string) []interface{} {
	for i, val := range s {
		s[i] = modifyData(val, replaceFunc)
	}
	return s
}

func DataToString(data interface{}, size int, workspaceId uint) string {
	indentedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}

	// dataWithEnvValues := utils.ReplaceByOperator(string(indentedJSON), workspaceId)
	dataWithEnvValues := string(indentedJSON)

	summariseProperties := func(s string) string {
		// s = utils.ReplaceByOperator(s, workspaceId)
		return utils.Truncate(s, size)
	}

	value, _ := replaceStringsInJSON(
		dataWithEnvValues,
		summariseProperties,
	)

	return value
}
