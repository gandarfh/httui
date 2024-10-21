package tree

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/pkg/styles"
	"github.com/gandarfh/httui/pkg/truncate"
)

var (
	bottomLeft = " └──"
)

type Styles struct {
	Title      lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
}

func defaultStyles() Styles {
	return Styles{
		Title:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff")),
		Selected:   lipgloss.NewStyle().Bold(true).Foreground(styles.DefaultTheme.PrimaryText),
		Unselected: lipgloss.NewStyle().Bold(false).Foreground(styles.DefaultTheme.FaintText),
	}
}

type Node[T any] struct {
	Value    string
	Data     T
	Children []Node[T]
	Expanded bool
}

type Model[T any] struct {
	Title     string
	KeyMap    KeyMap
	Styles    Styles
	width     int
	height    int
	cursor    int
	nodes     []Node[T]
	Paginator paginator.Model
}

func New[T any](nodes []Node[T], width int, height int) Model[T] {
	p := paginator.New()

	return Model[T]{
		KeyMap:    DefaultKeyMap(),
		Styles:    defaultStyles(),
		width:     width,
		height:    height,
		nodes:     nodes,
		Paginator: p,
	}
}

func (m Model[T]) Nodes() []Node[T] {
	return m.nodes
}

func (m *Model[T]) SetNodes(nodes []Node[T]) {
	m.nodes = nodes
	m.updatePagination()
}

func (m *Model[T]) NumberOfNodes() int {
	count := 0

	var countNodes func([]Node[T])
	countNodes = func(nodes []Node[T]) {
		for _, node := range nodes {
			count++
			if node.Expanded && node.Children != nil {
				countNodes(node.Children)
			}
		}
	}

	countNodes(m.nodes)
	return count
}

func (m *Model[T]) Index() int {
	log.Println(m.Paginator.Page*m.Paginator.PerPage, m.cursor)
	return m.Paginator.Page*m.Paginator.PerPage + m.cursor
}

func (m *Model[T]) Cursor() int {
	return m.cursor
}

func (m *Model[T]) SetCursorAndPage(cursor, page int) {
	m.cursor = cursor
	m.Paginator.Page = page
}

func (m *Model[T]) CurrentNode() *Node[T] {
	return m.GetNodeByIndex(m.Index())
}

func (m *Model[T]) GetNodeByIndex(index int) *Node[T] {
	count := 0

	var findNode func([]Node[T]) *Node[T]
	findNode = func(nodes []Node[T]) *Node[T] {
		for i := range nodes {
			if count == index {
				return &nodes[i]
			}
			count++

			if nodes[i].Expanded && len(nodes[i].Children) > 0 {
				if found := findNode(nodes[i].Children); found != nil {
					return found
				}
			}
		}
		return nil
	}

	return findNode(m.nodes)
}

func (m *Model[T]) ToggleExpand() {
	log.Println("ToggleExpand")
	count := 0
	var toggleExpandNode func([]Node[T]) bool

	toggleExpandNode = func(nodes []Node[T]) bool {
		cursor := m.Paginator.Page*m.Paginator.PerPage + m.cursor

		for i := range nodes {
			if count == cursor {
				nodes[i].Expanded = !nodes[i].Expanded
				return true
			}
			count++
			if nodes[i].Expanded && len(nodes[i].Children) > 0 {
				if toggleExpandNode(nodes[i].Children) {
					return true
				}
			}
		}
		return false
	}

	toggleExpandNode(m.nodes)
	m.updatePagination()
}

func (m *Model[T]) updatePagination() {
	totalItems := m.NumberOfNodes()
	m.Paginator.SetTotalPages(totalItems)
}

func (m *Model[T]) CursorUp() {
	if m.cursor > 0 {
		m.cursor--
	} else if m.Paginator.Page > 0 {
		m.Paginator.PrevPage()
		m.cursor = m.Paginator.PerPage - 1
	}
}

func (m *Model[T]) CursorDown() {
	itemsOnPage := m.Paginator.ItemsOnPage(m.NumberOfNodes())

	if m.cursor < itemsOnPage-1 {
		m.cursor++
	} else if !m.Paginator.OnLastPage() {
		m.Paginator.NextPage()
		m.cursor = 0
	}
}

func (m *Model[T]) CursorTop() {
	m.Paginator.Page = 0
	m.cursor = 0
}

func (m *Model[T]) CursorBottom() {
	log.Println(m.NumberOfNodes())
	m.Paginator.Page = m.Paginator.TotalPages - 1
	m.cursor = m.Paginator.ItemsOnPage(m.NumberOfNodes()) - 1
}

func (m *Model[T]) PrevPage() {
	m.Paginator.PrevPage()
}

func (m *Model[T]) NextPage() {
	m.Paginator.NextPage()
}

func (m Model[T]) Width() int {
	return m.width
}

func (m Model[T]) Height() int {
	return m.height
}

func (m *Model[T]) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model[T]) SetWidth(newWidth int) {
	m.SetSize(newWidth, m.height)
}

func (m *Model[T]) SetHeight(newHeight int) {
	m.SetSize(m.width, newHeight)
	m.Paginator.PerPage = m.Height() - 4
}

func (m Model[T]) Update(msg tea.Msg) (Model[T], tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Up):
			m.CursorUp()
		case key.Matches(msg, m.KeyMap.Down):
			m.CursorDown()

		case key.Matches(msg, m.KeyMap.Bottom):
			m.CursorBottom()
		case key.Matches(msg, m.KeyMap.Top):
			m.CursorTop()

		case key.Matches(msg, m.KeyMap.PrevPage):
			m.PrevPage()
		case key.Matches(msg, m.KeyMap.NextPage):
			m.NextPage()
		}
	}
	return m, nil
}

func (m Model[T]) View() string {
	if len(m.nodes) == 0 {
		return "No data"
	}

	lines := []string{
		m.Styles.Title.MaxWidth(m.width).Margin(1, 0).Render("   " + m.Title),
	}

	start, end := m.Paginator.GetSliceBounds(m.NumberOfNodes())
	renderedItems := m.renderTree(m.nodes, 0, &start, &end)

	for i, item := range renderedItems {
		valueWidth := m.width - m.width/6
		item = truncate.String(item, valueWidth)
		if i == m.cursor {
			lines = append(lines, m.Styles.Selected.MaxWidth(m.width).MarginBottom(1).Render(" > "+item))
		} else {
			lines = append(lines, m.Styles.Unselected.MaxWidth(m.width).MarginBottom(1).Render("   "+item))
		}
	}

	lines = append(lines, m.Paginator.View())

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m Model[T]) renderTree(remainingNodes []Node[T], indent int, start, end *int) []string {
	var nodesString []string

	var countNodes func([]Node[T], int)
	countNodes = func(nodes []Node[T], indent int) {
		for _, node := range nodes {
			if *start <= 0 && *end > 0 {
				var str string

				// If we aren't at the root, we add the arrow shape to the string
				if indent > 0 {
					shape := strings.Repeat(" ", (indent-1)*2) + bottomLeft + " "
					str += shape
				}

				valueStr := node.Value
				if node.Expanded {
					prefix := "[-] "

					if len(node.Children) == 0 {
						prefix = ""
					}

					valueStr = prefix + valueStr

				} else if len(node.Children) > 0 {
					valueStr = "[+] " + valueStr
				}
				str += valueStr

				nodesString = append(nodesString, str)
			}

			*start--
			*end--

			if node.Expanded && len(node.Children) > 0 {
				countNodes(node.Children, indent+1)
			}

			if *end <= 0 {
				return
			}
		}
	}

	countNodes(remainingNodes, indent)
	return nodesString
}

func MergeNodes[T any](originalNodes, updatedNodes []Node[T]) []Node[T] {
	var merge func([]Node[T], []Node[T]) []Node[T]
	merge = func(original, updated []Node[T]) []Node[T] {
		mergedNodes := make([]Node[T], len(original))

		for i := range original {
			mergedNodes[i] = original[i]

			for j := range updated {
				if original[i].Value == updated[j].Value {
					mergedNodes[i].Expanded = updated[j].Expanded

					if len(original[i].Children) > 0 && len(updated[j].Children) > 0 {
						mergedNodes[i].Children = merge(original[i].Children, updated[j].Children)
					}
					break
				}
			}
		}

		return mergedNodes
	}

	return merge(originalNodes, updatedNodes)
}
