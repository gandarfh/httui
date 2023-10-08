package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
)

type Page int

var (
	CurrPage      Page
	CurrEnv       repositories.Env
	CurrWorkspace repositories.Workspace
	CurrRequest   repositories.Request
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

type Command struct {
	Active   bool
	Value    string
	Category string
	Prefix   string
}

type CommandClose struct {
	Command
}

var command = Command{}

func SetCommand(value string) tea.Cmd {
	return func() tea.Msg {
		command.Value = value
		return command
	}
}

func OpenCommand(category, prefix string) tea.Cmd {
	return func() tea.Msg {
		command.Category = category
		command.Prefix = prefix
		command.Value = ""
		command.Active = true
		return command
	}
}

func CloseCommand() tea.Cmd {
	return func() tea.Msg {
		command.Active = false
		return CommandClose{command}
	}
}

func ClearCommand() tea.Cmd {
	return func() tea.Msg {
		command.Active = false
		command.Prefix = ""
		command.Value = ""
		command.Category = ""
		return command
	}
}

type Environment struct {
	Name string
}

func SetEnvironment(name string) tea.Cmd {
	return func() tea.Msg {
		return Environment{Name: name}
	}
}

type Loading struct {
	Msg   string
	Value bool
}

func SetLoading(loading bool, msg ...string) tea.Cmd {
	return func() tea.Msg {
		if len(msg) == 0 {
			return Loading{Msg: "", Value: loading}
		}

		return Loading{Msg: msg[0], Value: loading}
	}
}

func SetWorkspace(workspaceId uint) tea.Cmd {
	return func() tea.Msg {
		workspace_repo := repositories.NewWorkspace()
		CurrWorkspace, _ = workspace_repo.FindOne(workspaceId)

		return Loading{}
	}
}

type List struct {
	Requests []repositories.Request
}

var (
	ListOfWorkspaces []repositories.Workspace
	ListOfRequests   []repositories.Request
	ListOfEnvs       []repositories.Env
)

func ListRequests(parentId *uint) tea.Cmd {
	return func() tea.Msg {
		requests, _ := repositories.NewRequest().List(parentId, "")
		return List{Requests: requests}
	}
}

type Tab int

var (
	CurrTab = Tab_Requests
)

const (
	Tab_Requests Tab = iota
	Tab_Resources
)

func SetTab(tab Tab) tea.Cmd {
	return func() tea.Msg {
		CurrTab = tab
		return CurrTab
	}
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
