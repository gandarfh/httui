package common

import tea "github.com/charmbracelet/bubbletea"

type Page int

var (
	CurrPage Page
)

const (
	Page_Request Page = iota
	Page_Workspace
	Page_Resource
	Page_Env
)

func SetPage(page Page) tea.Cmd {
	return func() tea.Msg {
		CurrPage = Page(page)
		return page
	}
}

func SetNextPage() tea.Cmd {
	return func() tea.Msg {
		if CurrPage == 2 {
			CurrPage = 0
		} else {
			totalpages := 3
			CurrPage = Page(min(int(CurrPage)+1, totalpages-1))
		}
		return CurrPage
	}
}

func SetPrevPage() tea.Cmd {
	return func() tea.Msg {
		if CurrPage == 0 {
			CurrPage = 2
		} else {
			CurrPage = Page(max(int(CurrPage)-1, 0))
		}
		return CurrPage
	}
}
