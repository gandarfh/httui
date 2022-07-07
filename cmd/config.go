package cmd

import (
	"log"
	"os"

	"github.com/asaskevich/EventBus"
	"github.com/gandarfh/httui/cmd/model"
	"gopkg.in/yaml.v3"
)

var Bus = EventBus.New()

type Config struct {
	Default string      `yaml:"default"`
	Uris    []model.Uri `yaml:"dataSource"`
}

func (config *Config) GetUri(alias *string) *model.Uri {
	uri := model.Uri{}

	for _, item := range config.Uris {
		if item.Alias == *alias {

			uri = item
			break
		}
	}

	return &uri
}

func (config *Config) GetDefaultUri() *model.Uri {
	uri := *config.GetUri(&config.Default)

	return &uri
}

func (config *Config) SetDefaultUri() {}

func (config *Config) CreateUri() {}

func Connect() (*Config, error) {
	file, err := os.ReadFile("../config.yml")

	if err != nil {
		log.Fatal("Error when try read config.yml file:\n")
		return &Config{}, err
	}

	config := Config{}

	err = yaml.Unmarshal(file, &config)

	if err != nil {

		log.Fatal("Bad syntax in config.yml. \n", "Syntax must be like: `Config struct`")
		return &Config{}, err
	}

	return &config, nil
}
