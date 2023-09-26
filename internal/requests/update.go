package requests

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/repositories"
	"github.com/gandarfh/httui/pkg/common"
)

func (m Model) OpenRequest() Model {
	if len(common.ListOfRequests) == 0 {
		return m
	}

	index := m.List.Index()
	common.CurrRequest = common.ListOfRequests[index]

	if common.CurrRequest.Type == "group" {
		m.parentId = &common.CurrRequest.ID
		m.previousParentId = common.CurrRequest.ParentID
	}

	if common.CurrRequest.Type == "request" {
		m.parentId = common.CurrRequest.ParentID

		m.request_detail.Request = common.CurrRequest
		m.request_detail.Preview = fmt.Sprintf("%s - %s", common.CurrRequest.Method, common.CurrRequest.Endpoint)

		repositories.NewDefault().Update(&repositories.Default{
			RequestId: common.CurrRequest.ID,
		})
	}

	return m
}

func (m Model) BackRequest() Model {
	if m.parentId == nil {
		return m
	}

	if len(common.ListOfRequests) == 0 {
		m.parentId = m.previousParentId
		return m
	}

	group, _ := repositories.NewRequest().FindOne(*m.parentId)
	common.CurrRequest = *group

	m.parentId = common.CurrRequest.ParentID

	repositories.NewDefault().Update(&repositories.Default{
		RequestId: common.CurrRequest.ID,
	})

	return m
}

func (m Model) WindowSize(msg tea.WindowSizeMsg) Model {
	m.Height = msg.Height
	m.Width = msg.Width + 1

	m.List.SetHeight(m.Height/2 - 2)
	m.List.SetWidth(m.Width / 5)
	m.request_detail.Height = ((m.Height) - 9)
	m.request_detail.Width = m.Width - m.List.Width() + 1

	return m
}
