package data

import (
	"github/gandarfh/httui/config"
)

type Endpoints struct {
	Id      string `yaml:"id"`
	Path    string `yaml:"path"`
	Method  string `yaml:"method"`
	Headers string `yaml:"headers"`
	Body    string `yaml:"body"`
}

func (e Endpoints) Update(new Endpoints) error {
	connect := config.Connect[Data]()

	for _, item := range connect.Resources {
		for _, endpoint := range item.Endpoints {
			if endpoint.Id == e.Id {
				endpoint = new
				break
			}
		}
	}

	connect.Update(connect)

	return nil
}
