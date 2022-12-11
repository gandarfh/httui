package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	default_repo   *repositories.DefaultsRepo
	tags_list      list.Model
	resources_list list.Model
	filter         string
}

func New() Model {
	tags_repo, _ := repositories.NewTag()
	resources_repo, _ := repositories.NewResource()
	default_repo, _ := repositories.NewDefault()

	return Model{
		tags_repo:      tags_repo,
		resources_repo: resources_repo,
		default_repo:   default_repo,
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
		term := terminal.NewPreview(&msg.Response)

		return m, tea.Batch(common.SetLoading(false), term.OpenVim("Exec"))

	case terminal.Finish:
		switch msg.Category {
		case "Update":
			if common.CurrTab == common.Tab_Tags {
				index := m.tags_list.Index()
				common.CurrTag = common.ListOfTags[index]

				data := repositories.Tag{}
				msg.Preview.Execute(&data)
				m.tags_repo.Update(&common.CurrTag, &data)
			}

			if common.CurrTab == common.Tab_Resources {
				index := m.resources_list.Index()
				common.CurrResource = common.ListOfResources[index]

				msg.Preview.Execute(&common.CurrResource)
				m.resources_repo.Update(&common.CurrResource)

				return m, tea.Batch(common.ListResources(common.CurrTag.ID))
			}

		case "Create":
			if common.CurrTab == common.Tab_Tags {
				data := repositories.Tag{}
				msg.Preview.Execute(&data)
				// dont without name data
				if data.Name != "" {
					m.tags_repo.Create(&data)
				}
			}

			if common.CurrTab == common.Tab_Resources {
				data := repositories.Resource{}
				msg.Preview.Execute(&data)

				// dont without name data
				if data.Name != "" {
					m.resources_repo.Create(&data)
				}
			}

			defer msg.Preview.Close()
			if msg.Err != nil {
				return m, nil
			}
		}

	case common.Tab:
		common.CurrTab = msg
		common.ListOfResources = []repositories.Resource{}
		return m, tea.Batch(common.ClearResources())

	case common.List:
		common.ListOfTags = msg.Tags
		common.ListOfResources = msg.Resources

		m.tags_list.SetItems(m.TagsOfList())
		m.tags_list, cmd = m.tags_list.Update(msg)
		cmds = append(cmds, cmd)

		m.resources_list.SetItems(m.ResourcesOfList())
		m.resources_list, cmd = m.resources_list.Update(msg)
		cmds = append(cmds, cmd)

	case common.CommandClose:
		switch msg.Category {
		case "MOVE_TAG":
			m.ChangeTag(msg.Value)
			return m, common.ClearCommand()
		case "FILTER":
			m.filter = msg.Value
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height - 5
		m.width = msg.Width

		m.tags_list.SetHeight(msg.Height/2 - 2)
		m.resources_list.SetHeight(msg.Height/2 - 2)

	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			return m, tea.Batch(
				common.OpenCommand("FILTER"),
			)
		case "m":
			index := m.resources_list.Index()
			common.CurrResource = common.ListOfResources[index]

			return m, tea.Batch(
				common.OpenCommand("MOVE_TAG"),
				common.SetCommand(common.CurrResource.Tag.Name),
			)
		case "d":
			if common.CurrTab == common.Tab_Tags {
				index := m.tags_list.Index()
				m.tags_repo.Delete(common.ListOfTags[index].ID)
			}
			if common.CurrTab == common.Tab_Resources {
				index := m.resources_list.Index()
				m.resources_repo.Delete(common.ListOfResources[index].ID)
			}
		case "enter", "right", "l":
			return m, m.EnterResource()

		case "left", "h":
			return m, common.SetTab(common.Tab_Tags)

		case "e":
			if common.CurrTab == common.Tab_Resources {
				index := m.resources_list.Index()
				common.CurrResource = common.ListOfResources[index]

				return m, tea.Batch(common.SetLoading(true, "Loading..."), m.Exec())
			}

		case "esc":
			if common.CurrTab == common.Tab_Resources {
				m.resources_list.SetItems(nil)
				return m, common.SetTab(common.Tab_Tags)
			}

		case "r":
			var data interface{}
			if common.CurrTab == common.Tab_Tags {
				index := m.tags_list.Index()
				data = common.ListOfTags[index]
			}
			if common.CurrTab == common.Tab_Resources {
				index := m.resources_list.Index()
				data = common.ListOfResources[index]
			}

			term := terminal.NewPreview(&data)
			return m, tea.Batch(term.OpenVim("Update"))
		case "c":
			var data interface{}
			if common.CurrTab == common.Tab_Tags {
				data = repositories.Tag{WorkspaceId: common.CurrWorkspace.ID}
			}
			if common.CurrTab == common.Tab_Resources {
				data = repositories.Resource{TagId: common.CurrTag.ID}
			}

			term := terminal.NewPreview(&data)
			return m, tea.Batch(term.OpenVim("Create"))
		}
	}

	if common.CurrTab == common.Tab_Tags {
		common.ListOfTags, _ = m.tags_repo.List(common.CurrWorkspace.ID)
		m.tags_list.SetItems(m.TagsOfList())
		m.tags_list, cmd = m.tags_list.Update(msg)

		cmds = append(cmds, cmd)
	}

	m.resources_list.SetItems(m.ResourcesOfList())
	if common.CurrTab == common.Tab_Resources {
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

	m.resources_list.Title = fmt.Sprintf("%s -> %s", common.CurrWorkspace.Name, common.CurrTag.Name)

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
			value := utils.ReplaceByEnv(v)
			res.Header(k, value)
			headers[k] = value
		}

		params := []map[string]string{}
		json.Unmarshal(resource.QueryParams, &params)
		for _, item := range params {
			for k, v := range item {
				value := utils.ReplaceByEnv(v)
				res.Params(k, value)
				params = append(params, map[string]string{k: value})
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

func (m Model) ChangeTag(newtag string) error {
	tag, err := m.tags_repo.FindOneByname(newtag, common.CurrTag.WorkspaceId)
	if err != nil {
		log.Fatal(err)
	}

	index := m.resources_list.Index()
	common.CurrResource = common.ListOfResources[index]

	common.CurrResource.TagId = tag.ID
	m.resources_repo.Update(&common.CurrResource)

	return nil
}

func (m Model) EnterResource() tea.Cmd {
	if common.CurrTab == common.Tab_Tags && len(common.ListOfTags) > 0 {
		index := m.tags_list.Index()
		common.CurrTag = common.ListOfTags[index]

		data := repositories.Default{
			TagId: common.CurrTag.ID,
		}

		m.default_repo.Update(&data)

		if len(common.CurrTag.Resources) > 0 {
			return common.SetTab(common.Tab_Resources)
		}
	}

	return nil
}

func Curl(req *http.Request) string {
	command, _ := http2curl.GetCurlCommand(req)
	return command.String()
}
