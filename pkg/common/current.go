package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/maid-san/internal/repositories"
)

type Page = int

func SetPage(page int) tea.Cmd {
	return func() tea.Msg {
		CurrPage = page
		return page
	}
}

const (
	Page_Workspace Page = 0
	Page_Resource  Page = 1
	Page_Env       Page = 2
)

var (
	CurrPage      Page
	CurrWorkspace repositories.Workspace
	CurrTag       repositories.Tag
)
