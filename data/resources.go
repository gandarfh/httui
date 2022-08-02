package data

import (
	"github.com/google/uuid"
	"github/gandarfh/httui/config"
)

type Resources struct {
	Id        string      `yaml:"id"`
	Alias     string      `yaml:"alias"`
	Endpoints []Endpoints `yaml:"endpoints"`
}

func (r Resources) New() error {
	connect := config.Connect[Data]()

	new := Resources{
		Id:        uuid.NewString(),
		Alias:     r.Alias,
		Endpoints: r.Endpoints,
	}

	connect.Resources = append(connect.Resources, new)
	connect.Update(connect)

	return nil
}

func (r Resources) Update(id string) error {
	connect := config.Connect[Data]()

	new := Resources{
		Id:        id,
		Alias:     r.Alias,
		Endpoints: r.Endpoints,
	}

	for _, item := range connect.Resources {
		if item.Id == id {
			item = new
			break
		}
	}

	connect.Update(connect)

	return nil
}

func (r Resources) NewEndpoint(endpoint Endpoints) error {
	connect := config.Connect[Data]()

	new := Endpoints{
		Id:      uuid.NewString(),
		Path:    endpoint.Path,
		Method:  endpoint.Method,
		Headers: endpoint.Headers,
		Body:    endpoint.Body,
	}

	for _, item := range connect.Resources {
		if item.Id == r.Id {
			item.Endpoints = append(item.Endpoints, new)
		}
	}

	connect.Update(connect)

	return nil
}
