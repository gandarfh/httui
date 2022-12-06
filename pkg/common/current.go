package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/maid-san/internal/repositories"
)

type List struct {
	Tags      []repositories.Tag
	Resources []repositories.Resource
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

type Page = int

func SetPage(page int) tea.Cmd {
	return func() tea.Msg {
		CurrPage = page
		return page
	}
}

const (
	Page_Workspace Page = iota
	Page_Resource
	Page_Env
)

type ResourceTab struct {
	Active int
}

const (
	Tab_Tags = iota
	Tab_Resources
)

func SetResourceTab(tab int) tea.Cmd {
	return func() tea.Msg {
		CurrResoruceTab.Active = tab
		return CurrResoruceTab
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

var (
	CurrPage        Page
	CurrResoruceTab = ResourceTab{Active: Tab_Tags}

	CurrEnv       repositories.Env
	CurrWorkspace repositories.Workspace
	CurrTag       repositories.Tag
	CurrResource  repositories.Resource

	ListOfWorkspaces []repositories.Workspace
	ListOfTags       []repositories.Tag
	ListOfResources  []repositories.Resource
	ListOfEnvs       []repositories.Env
)
