package resources

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/maid-san/internal/repositories"
	"github.com/gandarfh/maid-san/pkg/client"
	"github.com/gandarfh/maid-san/pkg/common"
	"github.com/gandarfh/maid-san/pkg/terminal"
	"github.com/gandarfh/maid-san/pkg/utils"
	"moul.io/http2curl"
)

var (
	divider = lipgloss.NewStyle().MarginLeft(1).Border(lipgloss.NormalBorder(), false, true, false, false)
)

type Model struct {
	width          int
	height         int
	tags_repo      *repositories.TagsRepo
	resources_repo *repositories.ResourcesRepo
	tags_list      list.Model
	resources_list list.Model
}

func New() Model {
	tags_repo, _ := repositories.NewTag()
	resources_repo, _ := repositories.NewResource()

	return Model{
		tags_repo:      tags_repo,
		resources_repo: resources_repo,
		tags_list:      NewTagList(),
		resources_list: NewResourceList(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case Result:
		// m.Loading = false
		term := terminal.NewPreview(&msg.Response)

		return m, tea.Batch(term.OpenVim("Exec"))

	case terminal.Finish:
		switch msg.Category {
		case "Update":
			if common.CurrResoruceTab.Active == common.Tab_Tags {
				index := m.tags_list.Index()
				common.CurrTag = common.ListOfTags[index]

				data := repositories.Tag{}
				msg.Preview.Execute(&data)
				m.tags_repo.Update(&common.CurrTag, &data)
			}

			if common.CurrResoruceTab.Active == common.Tab_Resources {
				index := m.resources_list.Index()
				common.CurrResource = common.ListOfResources[index]

				data := repositories.Resource{}
				msg.Preview.Execute(&data)
				m.resources_repo.Update(&common.CurrResource, &data)

				return m, tea.Batch(common.ListResources(common.CurrTag.ID))
			}

		case "Create":
			if common.CurrResoruceTab.Active == common.Tab_Tags {
				data := repositories.Tag{}
				msg.Preview.Execute(&data)
				m.tags_repo.Create(&data)
			}

			if common.CurrResoruceTab.Active == common.Tab_Resources {
				data := repositories.Resource{}
				msg.Preview.Execute(&data)
				m.resources_repo.Create(&data)
			}

			defer msg.Preview.Close()
			if msg.Err != nil {
				return m, nil
			}
		}

	case common.ResourceTab:
		common.CurrResoruceTab = msg
		common.ListOfResources = []repositories.Resource{}
		return m, tea.Batch(common.ClearResources())
		// return m, nil

	case common.List:
		common.ListOfTags = msg.Tags
		common.ListOfResources = msg.Resources

		m.tags_list.SetItems(m.TagsOfList())
		m.tags_list, cmd = m.tags_list.Update(msg)
		cmds = append(cmds, cmd)

		m.resources_list.SetItems(m.ResourcesOfList())
		m.resources_list, cmd = m.resources_list.Update(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.height = msg.Height - (msg.Height / 3)
		m.width = msg.Width

		m.tags_list.SetHeight(14)
		m.resources_list.SetHeight(14)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "d":
			if common.CurrResoruceTab.Active == common.Tab_Tags {
				index := m.tags_list.Index()
				m.tags_repo.Delete(common.ListOfTags[index].ID)
			}
			if common.CurrResoruceTab.Active == common.Tab_Resources {
				index := m.resources_list.Index()
				m.resources_repo.Delete(common.ListOfResources[index].ID)
			}
		case "tab":
			return m, tea.Batch(common.SetResourceTab(common.Tab_Resources))
		case "enter":
			index := m.tags_list.Index()
			common.CurrTag = common.ListOfTags[index]

			m.resources_list.SetItems(m.ResourcesOfList())
			m.resources_list, cmd = m.resources_list.Update(msg)
			return m, tea.Batch(common.SetResourceTab(common.Tab_Resources))

		case "e":
			// m.Loading = true
			index := m.resources_list.Index()
			common.CurrResource = common.ListOfResources[index]

			return m, tea.Batch(m.Exec())

		case "esc", "shift+tab":

			m.resources_list.SetItems(nil)
			return m, tea.Batch(common.SetResourceTab(common.Tab_Tags))

		case "r":
			var data interface{}
			if common.CurrResoruceTab.Active == common.Tab_Tags {
				index := m.tags_list.Index()
				data = common.ListOfTags[index]
			}
			if common.CurrResoruceTab.Active == common.Tab_Resources {
				index := m.resources_list.Index()
				data = common.ListOfResources[index]
			}

			term := terminal.NewPreview(&data)
			return m, tea.Batch(term.OpenVim("Update"))
		case "c":
			var data interface{}
			if common.CurrResoruceTab.Active == common.Tab_Tags {
				data = repositories.Tag{WorkspaceId: common.CurrWorkspace.ID}
			}
			if common.CurrResoruceTab.Active == common.Tab_Resources {
				data = repositories.Resource{TagId: common.CurrTag.ID}
			}

			term := terminal.NewPreview(&data)
			return m, tea.Batch(term.OpenVim("Create"))
		}
	}

	if common.CurrResoruceTab.Active == common.Tab_Tags {
		m.tags_list.SetItems(m.TagsOfList())
		m.tags_list, cmd = m.tags_list.Update(msg)

		cmds = append(cmds, cmd)
	}

	if common.CurrResoruceTab.Active == common.Tab_Resources {
		common.ListOfTags, _ = m.tags_repo.List(common.CurrWorkspace.ID)
		m.resources_list.SetItems(m.ResourcesOfList())
		m.resources_list, cmd = m.resources_list.Update(msg)

		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if common.CurrWorkspace.Name != "" {
		common.ListOfTags, _ = m.tags_repo.List(common.CurrWorkspace.ID)
		m.tags_list.Title = common.CurrWorkspace.Name
	} else {
		m.tags_list.Title = "Select some workspace!"
	}

	m.resources_list.Title = common.CurrTag.Name

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.tags_list.View(),
		divider.Height(m.height).String(),
		m.resources_list.View(),
	)
}

type Result struct {
	Err      error
	Response any
	Loading  bool
}

func (m Model) Exec() tea.Cmd {
	return func() tea.Msg {
		workspace := common.CurrWorkspace
		resource := common.CurrResource

		url := utils.ReplaceByEnv(workspace.Host) + utils.ReplaceByEnv(resource.Endpoint)

		res := client.Request(url, strings.ToUpper(resource.Method))

		rawbody, _ := resource.Body.MarshalJSON()
		bodystring := utils.ReplaceByEnv(string(rawbody))

		var body any
		if err := json.Unmarshal([]byte(bodystring), &body); err != nil {
			panic(err)
		}

		if _, ok := body.(map[string]any); ok {
			res.Body([]byte(bodystring))
		} else {
			res.Body(nil)
		}

		headers := map[string]string{}
		json.Unmarshal(resource.Headers, &headers)

		for k, v := range headers {
			res.Header(k, utils.ReplaceByEnv(v))
		}

		params := []map[string]string{}
		json.Unmarshal(resource.Headers, &headers)
		for _, item := range params {
			for k, v := range item {
				res.Params(k, utils.ReplaceByEnv(v))
			}
		}

		data, _ := res.Exec()

		var response any
		readbody, _ := ioutil.ReadAll(data.Body)
		json.Unmarshal(readbody, &response)
		result := struct {
			Url      string              `json:"url"`
			Method   string              `json:"method"`
			Status   string              `json:"status"`
			Params   []map[string]string `json:"params"`
			Headers  map[string]string   `json:"headers"`
			Body     any                 `json:"body"`
			Response any                 `json:"response"`
			Curl     string              `json:"curl"`
		}{
			Url:      url,
			Method:   resource.Method,
			Status:   data.Status,
			Params:   params,
			Headers:  headers,
			Body:     resource.Body,
			Response: response,
			Curl:     Curl(data.Request),
		}

		return Result{
			Response: result,
		}
	}
}

func Curl(req *http.Request) string {
	command, _ := http2curl.GetCurlCommand(req)
	return command.String()
}
