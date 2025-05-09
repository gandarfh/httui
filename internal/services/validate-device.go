package services

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type OrganizationRole struct {
	ID   string `json:"_id,omitempty"`
	Role string `json:"role,omitempty"`
}

type User struct {
	ID            string             `json:"_id,omitempty"`
	CreatedAt     time.Time          `json:"createdAt,omitempty"`
	UpdatedAt     time.Time          `json:"updatedAt,omitempty"`
	DeletedAt     *time.Time         `json:"deletedAt,omitempty"`
	Organizations []OrganizationRole `json:"organizations,omitempty"`
	ProviderId    string             `json:"providerId,omitempty"`
	Image         string             `json:"image,omitempty"`
	Email         string             `json:"email,omitempty"`
	FirstName     string             `json:"firstName,omitempty"`
	LastName      string             `json:"lastName,omitempty"`
}

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type ValidateDeviceResponse struct {
	User   `json:",inline"`
	Tokens *Tokens `json:"tokens"`
}

func PollingValidate(deviceId, accessToken string) tea.Cmd {
	return tea.Tick(5*time.Second, func(time.Time) tea.Msg {
		var response = ValidateDeviceResponse{}

		client, err := HttuiApiDatasource.
			Auth("Authorization", fmt.Sprint("Bearer ", accessToken)).
			Post("/auth/validate/device/" + deviceId)
		if err != nil {
			return ValidateDeviceResponse{}
		}

		if err := client.Decode(&response); err != nil {
			return ValidateDeviceResponse{}
		}

		return response
	})
}
