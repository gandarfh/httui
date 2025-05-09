package services

import (
	"log"
	"time"
)

type PlanPeriodType string

const (
	PlanPeriod_Monthly PlanPeriodType = "monthly"
	PlanPeriod_Weekly  PlanPeriodType = "weekly"
	PlanPeriod_Daily   PlanPeriodType = "daily"
)

type Period struct {
	Number int            `json:"number,omitempty"`
	Type   PlanPeriodType `json:"type,omitempty"`
}

type Price struct {
	Precision int    `json:"precision,omitempty"`
	Amount    int64  `json:"amount,omitempty"`
	Currency  string `json:"currency,omitempty"`
}

type Plan struct {
	ID        string     `json:"_id,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	Name      string     `json:"name,omitempty"`
	Price     Price      `json:"price,omitempty"`
	Period    Period     `json:"period,omitempty"`
}

type BillingStatus string

const (
	BillingStatus_Succeeded BillingStatus = "succeeded"
	BillingStatus_Cancelled BillingStatus = "cancelled"
	BillingStatus_Pending   BillingStatus = "pending"
)

type Billing struct {
	ID             string        `json:"_id,omitempty"`
	CreatedAt      time.Time     `json:"createdAt,omitempty"`
	UpdatedAt      time.Time     `json:"updatedAt,omitempty"`
	DeletedAt      *time.Time    `json:"deletedAt,omitempty"`
	OrganizationId string        `json:"organizationId,omitempty"`
	Status         BillingStatus `json:"status,omitempty"`
	Price          Price         `json:"price,omitempty"`
	Period         Plan          `json:"plan,omitempty"`
	ExternalId     string        `json:"externalId,omitempty"`
}

type Organization struct {
	ID        string     `json:"_id,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	Name      string     `json:"name,omitempty"`
	Plan      *Plan      `json:"plan,omitempty"`
	Billings  []Billing  `json:"billings,omitempty"`
}

type OrganizationResponse struct {
	*Organization `json:",inline"`
	MqttTopic     string `json:"mqttTopic"`
	MqttEndpoint  string `json:"mqttEndpoint"`
	MqttCertPEM   string `json:"mqttCertPEM"`
	MqttCertKey   string `json:"mqttCertKey"`
	MqttCertCA    string `json:"mqttCertCA"`
}

func OrganizationShow() (*OrganizationResponse, error) {
	client, err := HttuiApiDatasource.Get("organizations")
	if err != nil {
    log.Println("OrganizationShow:",err.Error())
		return nil, err
	}

	var response OrganizationResponse
	if err := client.Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
