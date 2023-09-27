package requests

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/utils"
)

var (
	styleConfig = glamour.DarkStyleConfig
)

type ModelDetail struct {
	Width   int
	Height  int
	Preview string
	Request repositories.Request
}

func NewDetail() ModelDetail {
	styleConfig.CodeBlock.Chroma.Error.BackgroundColor = nil

	return ModelDetail{}
}

func (m ModelDetail) Init() tea.Cmd {
	return nil
}

func (m ModelDetail) Update(msg tea.Msg) (ModelDetail, tea.Cmd) {
	var (
		_    tea.Cmd
		cmds []tea.Cmd
	)

	return m, tea.Batch(cmds...)
}

func DataToString(jsonString string, size int) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return ""
	}

	// Em seguida, serialize a estrutura de dados de volta em JSON com indentação
	indentedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}

	return utils.Truncate(utils.ReplaceByEnv(string(indentedJSON)), size)
}

func (m ModelDetail) View() string {
	title_row := lipgloss.NewStyle().Width(m.Width).Border(lipgloss.RoundedBorder()).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Width(m.Width-m.Width/3).Bold(true).Render(" Name: "+m.Request.Name),
		lipgloss.NewStyle().Width(m.Width/6).String(),
		lipgloss.NewStyle().Width(m.Width/6).Bold(true).Render(fmt.Sprint(" ID: ", m.Request.ID)),
	))

	preview_row := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(
		lipgloss.Place(
			m.Width,
			1,
			lipgloss.Left,
			lipgloss.Top,
			lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Preview: ", utils.ReplaceByEnv(utils.Truncate(m.Preview, 100)))),
		))

	bodyrenderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(styleConfig),
		glamour.WithWordWrap(m.Width-m.Width/3),
	)

	paramrenderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(styleConfig),
		glamour.WithWordWrap(m.Width/3-10),
	)

	rawbody, _ := m.Request.Body.MarshalJSON()
	body, _ := bodyrenderer.Render("```json\n" + DataToString(string(rawbody), 400) + "\n```")

	body_box := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Body:")),
		body,
	)

	rawparams, _ := m.Request.QueryParams.MarshalJSON()
	query, _ := paramrenderer.Render("```json\n" + DataToString(fmt.Sprintf(`{"items": %s}`, string(rawparams)), 80) + "\n```")

	query_box := lipgloss.NewStyle().Height(m.Height / 2).Render(lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Params:")),
		query,
	))

	rawheader, _ := m.Request.Headers.MarshalJSON()
	header, _ := paramrenderer.Render("```json\n" + DataToString(fmt.Sprintf(`{"items": %s}`, string(rawheader)), 80) + "\n```")

	header_box := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Headers:")),
		header,
	)

	params_box := lipgloss.
		NewStyle().
		Height(m.Height).
		Width(m.Width / 3).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			query_box,
			lipgloss.NewStyle().Width(m.Width/3).Border(lipgloss.NormalBorder(), true, false, false, false).String(),
			header_box,
		))

	content_row := lipgloss.NewStyle().Width(m.Width).Height(m.Height).Border(lipgloss.RoundedBorder()).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		body_box,
		lipgloss.NewStyle().Height(m.Height).Width(1).Border(lipgloss.NormalBorder(), false, true, false, false).String(),
		params_box,
	))

	container := lipgloss.NewStyle().Padding(0, 1).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		title_row,
		preview_row,
		content_row,
	))

	return container
}

func (m ModelDetail) Help() string {
	return ""
}
