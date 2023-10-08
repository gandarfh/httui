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

func DataToString(data interface{}, size int) string {
	indentedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}

	return utils.Truncate(utils.ReplaceByOperator(string(indentedJSON)), size)
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
			lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Preview: ", utils.ReplaceByOperator(utils.Truncate(m.Preview, 100)))),
		))

	bodyrenderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(styleConfig),
		glamour.WithWordWrap(m.Width-m.Width/3),
	)

	paramrenderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(styleConfig),
		glamour.WithWordWrap(m.Width/3-10),
	)

	rawbody := m.Request.Body.Data()
	body, _ := bodyrenderer.Render("```json\n" + DataToString(rawbody, 400) + "\n```")

	body_box := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Body:")),
		body,
	)

	rawparams := utils.GetAllParentsParams(m.Request.ParentID, m.Request.QueryParams.Data())
	query, _ := paramrenderer.Render("```json\n" + DataToString(rawparams, 80) + "\n  ```")

	query_box := lipgloss.NewStyle().Height(m.Height / 2).Render(lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprint(" Params:")),
		query,
	))

	rawheader := utils.GetAllParentsHeaders(m.Request.ParentID, m.Request.Headers.Data())
	header, _ := paramrenderer.Render("```json\n" + DataToString(rawheader, 80) + "\n  ```")

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
