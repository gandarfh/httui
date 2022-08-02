package data

import (
	"github/gandarfh/httui/config"
	"github/gandarfh/httui/utils"
	"os"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Default   string      `yaml:"default"`
	Resources []Resources `yaml:"dataSource"`
}

func (d Data) Update(content Data) (Data, error) {
	out, err := yaml.Marshal(&content)

	if err != nil {
		utils.MsgError(err)
		return Data{}, err
	}

	err = os.WriteFile("config.yml", out, 0666)

	if err != nil {
		utils.MsgError(err)
		return Data{}, err
	}

	return content, nil
}

func (d Data) GetAllEndpoints() []Endpoints {
	connect := config.Connect[Data]()

	endpoints := []Endpoints{}

	for _, item := range connect.Resources {
		endpoints = append(endpoints, item.Endpoints...)
	}

	return endpoints
}

func (d Data) GetAllResources() []Resources {
	connect := config.Connect[Data]()

	resources := []Resources{}

	for _, item := range connect.Resources {
		resources = append(resources, item)
	}

	return resources
}

func (d Data) GetResource(alias string) Resources {
	connect := config.Connect[Data]()

	resources := Resources{}

	for _, item := range connect.Resources {
		if item.Alias == alias {

			resources = item
			break
		}
	}

	return resources
}

func (d Data) GetDefaultResources() Resources {
	connect := config.Connect[Data]()

	resources := d.GetResource(connect.Default)

	return resources
}
