package requests

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/pkg/tree/v2"
)

func LoadDefault() tea.Msg {
	config, _ := offline.NewDefault().First()
	return *config
}

func LoadWorspace() tea.Msg {
	config, _ := offline.NewDefault().First()
	workspace, _ := offline.NewWorkspace().FindOne(config.WorkspaceId)
	log.Println(workspace.ID, workspace.Name)

	return workspace
}

type RequestsData struct {
	List        []offline.Request
	Current     offline.Request
	RequestTree []tree.Node[offline.Request]
	ParentID    *uint
	Cursor      int
	Page        int
}

func LoadRequests() tea.Msg {
	config, _ := offline.NewDefault().First()
	request, _ := offline.NewRequest().FindOne(config.RequestId)
	requests, _ := offline.NewRequest().List(request.ParentID, "")

	return RequestsData{
		RequestTree: config.RequestTree.Data(),
		Cursor:      *config.Cursor,
		Page:        *config.Page,
		Current:     *request,
		List:        requests,
		ParentID:    request.ParentID,
	}
}

func LoadRequestsByParentId(parentId *uint) tea.Cmd {
	return func() tea.Msg {
		config, _ := offline.NewDefault().First()
		request, _ := offline.NewRequest().FindOne(config.RequestId)
		requests, _ := offline.NewRequest().List(parentId, "")

		return RequestsData{
			List:        requests,
			ParentID:    parentId,
			RequestTree: config.RequestTree.Data(),
			Cursor:      *config.Cursor,
			Page:        *config.Page,
			Current:     *request,
		}
	}
}

func LoadRequestsByFilter(filter string) tea.Cmd {
	return func() tea.Msg {
		config, _ := offline.NewDefault().First()
		request, _ := offline.NewRequest().FindOne(config.RequestId)
		requests, _ := offline.NewRequest().List(nil, filter)

		return RequestsData{
			RequestTree: config.RequestTree.Data(),
			List:        requests,
			Current:     *request,
			ParentID:    nil,
			Cursor:      0,
			Page:        0,
		}
	}
}
