package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/asaskevich/EventBus"
	"github.com/gandarfh/httui/cmd/model"
	"gopkg.in/yaml.v3"
)

var Bus = EventBus.New()

type Config struct {
	Default string       `yaml:"default"`
	Uris    *[]model.Uri `yaml:"dataSource"`
}

func (config *Config) GetUri(alias *string) *model.Uri {
	uri := model.Uri{}

	for _, item := range *config.Uris {
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

func (config *Config) SetDefaultUri(uri *string) {
	newConfig := Config{
		Uris:    config.Uris,
		Default: *uri,
	}

	err := updateConfig(&newConfig)

	if err != nil {
		log.Panicln(err)
	}
}

func (config *Config) CreateUri(uri *model.Uri) error {
	*config.Uris = append(*config.Uris, *uri)

	if err := updateConfig(config); err != nil {
		return err
	}

	return nil
}

func (config *Config) CreateEndpoint(uri *model.Uri, in *model.Endpoint) error {

	for _, item := range *config.Uris {
		if item.Alias == uri.Alias {

			*item.Endpoints = append(*item.Endpoints, *in)
			break
		}
	}

	fmt.Print(config.GetUri(&uri.Alias))

	if err := updateConfig(config); err != nil {
		return err
	}

	return nil
}

func Connect() (*Config, error) {
	file, err := os.ReadFile("config.yml")

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error when try read config.yml file:\n")
		return &Config{}, err
	}

	config := Config{}

	err = yaml.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err)
		log.Fatal("Bad syntax in config.yml. \n", "Syntax must be like: `Config struct`")
		return &Config{}, err
	}

	return &config, nil
}

func updateConfig(in *Config) error {
	out, err := yaml.Marshal(&in)

	if err != nil {

		return err
	}

	err = os.WriteFile("config.yml", out, 0666)

	if err != nil {

		return err
	}

	return nil
}
