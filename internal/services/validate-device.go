package services

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gandarfh/httui/pkg/client"
)

type Profile struct {
	ID             string     `json:"_id,omitempty"`
	CreatedAt      time.Time  `json:"createdAt,omitempty"`
	UpdatedAt      time.Time  `json:"updatedAt,omitempty"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty"`
	OrganizationId string     `json:"organizationId,omitempty"`
	ProviderId     string     `json:"providerId,omitempty"`
	Image          string     `json:"image,omitempty"`
	Email          string     `json:"email,omitempty"`
	FirstName      string     `json:"firstName,omitempty"`
	LastName       string     `json:"lastName,omitempty"`
}

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type ValidateDeviceResponse struct {
	Profile `json:",inline"`
	Tokens  *Tokens `json:"tokens"`
}

func PollingValidate(accessToken string) tea.Cmd {
	return tea.Tick(5*time.Second, func(time.Time) tea.Msg {
		var response = ValidateDeviceResponse{}

		api, err := client.Post("http://localhost:5000/auth/validate/device").
			Header("Authorization", fmt.Sprint("Bearer ", accessToken)).
			Exec()

		if err != nil {
			return ValidateDeviceResponse{}
		}

		readbody, err := io.ReadAll(api.Body)
		if err != nil {
			return ValidateDeviceResponse{}
		}

		if err := json.Unmarshal(readbody, &response); err != nil {
			return ValidateDeviceResponse{}
		}

		return response
	})
}
