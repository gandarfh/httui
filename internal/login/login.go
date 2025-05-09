package login

import (
	"fmt"
	"os/user"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/internal/config"
	"github.com/gandarfh/httui/internal/services"
	"github.com/gandarfh/httui/pkg/browser"
)

type Model struct {
	Width       int
	Height      int
	keys        KeyMap
	url         string
	success     bool
	accessToken string
	deviceID    string
}

func New() Model {
	defaultConfig := config.ConfigParser{}
	config.Config = defaultConfig.GetDefaultConfig()

	return Model{
		keys: keys,
	}
}

func (m Model) Init() tea.Cmd {
	name, _ := user.Current()
	device := services.Device{
		Name: name.Name,
	}

	return device.Create
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case services.ValidateDeviceResponse:
		if msg.Tokens != nil {
			m.success = true
			config.Config.Settings.Token = msg.Tokens.Access
			config.UpdateConfig(config.Config)
			return m, tea.Quit
		}

		return m, services.PollingValidate(m.deviceID, m.accessToken)

	case services.DeviceResponse:
		m.url = msg.Url
		config.Config.Settings.DeviceID = msg.ID
		m.deviceID = msg.ID
		m.accessToken = msg.Tokens.Access

		return m, services.PollingValidate(msg.ID, msg.Tokens.Access)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.OpenPage):
			cmd = browser.OpenPage(m.url)
			return m, cmd
		}
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.url == "" {
		return "loading..."
	}

	if m.success {
		return "Device connected with success!"
	}

	return fmt.Sprintf("Press Enter to open the browser or visit %s", m.url)
}
