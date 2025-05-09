package services

import (
	"encoding/json"
	"io"
	"net/http"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type DeviceStatus string

const (
	DeviceStatus_Pending  DeviceStatus = "Pending"
	DeviceStatus_Active   DeviceStatus = "Active"
	DeviceStatus_Inactive DeviceStatus = "Inactive"
	DeviceStatus_Revoked  DeviceStatus = "Revoked"
)

type Device struct {
	ID             string       `json:"_id,omitempty"`
	CreatedAt      time.Time    `json:"createdAt,omitempty"`
	UpdatedAt      time.Time    `json:"updatedAt,omitempty"`
	DeletedAt      *time.Time   `json:"deletedAt,omitempty"`
	ProfileId      string       `json:"profileId,omitempty"`
	Name           string       `json:"name,omitempty"`
	Url            string       `json:"url,omitempty"`
	Platform       string       `json:"platform,omitempty"`
	IPAddress      string       `json:"ipaddress,omitempty"`
	Status         DeviceStatus `json:"status,omitempty"`
	LastAccessedAt time.Time    `json:"lastAccessedAt,omitempty"`
	RevocationAt   *time.Time   `json:"revocationAt,omitempty"`
}

type DeviceResponse struct {
	Device `json:",inline"`
	Tokens *Tokens `json:"tokens"`
}

func (d Device) Create() tea.Msg {
	body := map[string]string{
		"name":      d.Name,
		"ipaddress": GetPublicIP(),
		"platform":  runtime.GOOS,
	}

	data, _ := json.Marshal(body)

	client, _ := HttuiApiDatasource.Body(data).Post("/devices")
	response := DeviceResponse{}
	client.Decode(&response)

	return response
}

func GetPublicIP() string {
	response, err := http.Get("https://api.ipify.org")
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	ip, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	return string(ip)
}
