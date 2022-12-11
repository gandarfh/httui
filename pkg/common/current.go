package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/maid-san/internal/repositories"
)

type Page = int

var (
	CurrPage      Page
	CurrEnv       repositories.Env
	CurrWorkspace repositories.Workspace
	CurrTag       repositories.Tag
	CurrResource  repositories.Resource
)

const (
	Page_Workspace Page = iota
	Page_Resource
	Page_Env
)

func SetPage(page int) tea.Cmd {
	return func() tea.Msg {
		CurrPage = page
		return page
	}
}

type Command struct {
	Active bool
	Value  string
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

func OpenCommand() tea.Cmd {
	return func() tea.Msg {
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
		command.Value = ""
		return command
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
		workspace_repo, _ := repositories.NewWorkspace()
		CurrWorkspace, _ = workspace_repo.FindOne(workspaceId)

		return Loading{}
	}
}

func SetTag(tagId uint) tea.Cmd {
	return func() tea.Msg {
		tags_repo, _ := repositories.NewTag()
		CurrTag, _ = tags_repo.FindOne(tagId)

		return Loading{}
	}
}

type List struct {
	Tags      []repositories.Tag
	Resources []repositories.Resource
}

var (
	ListOfWorkspaces []repositories.Workspace
	ListOfTags       []repositories.Tag
	ListOfResources  []repositories.Resource
	ListOfEnvs       []repositories.Env
)

func ListTags(workspaceId uint) tea.Cmd {
	return func() tea.Msg {
		tags_repo, _ := repositories.NewTag()
		ListOfTags, _ = tags_repo.List(workspaceId)

		return List{Tags: ListOfTags}
	}
}

func ListResources(tagId uint) tea.Cmd {
	return func() tea.Msg {
		resource_repo, _ := repositories.NewResource()
		ListOfResources, _ = resource_repo.List(tagId)
		return List{Resources: ListOfResources}
	}
}

func ClearResources() tea.Cmd {
	return func() tea.Msg {
		ListOfResources = []repositories.Resource{}
		return List{Resources: []repositories.Resource{}}
	}
}

type Tab int

var (
	CurrTab = Tab_Tags
)

const (
	Tab_Tags Tab = iota
	Tab_Resources
)

func SetTab(tab Tab) tea.Cmd {
	return func() tea.Msg {
		CurrTab = tab
		return CurrTab
	}
}
