package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	repositoryenvs "github.com/gandarfh/maid-san/internal/commands/envs/repository"
	"github.com/gandarfh/maid-san/internal/commands/generator/dtos"
	resourcesdtos "github.com/gandarfh/maid-san/internal/commands/resources/dtos"
	workspacedtos "github.com/gandarfh/maid-san/internal/commands/workspace/dtos"
	repositorywork "github.com/gandarfh/maid-san/internal/commands/workspace/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/validate"
)

type parameters struct {
	Name   string `json:"name"`
	In     string `json:"in"` // "body" | "query" | "header"
	Schema struct {
		Ref string `json:"$ref"`
	} `json:"schema"`
}

type SwaggerData struct {
	Swagger string `json:"swagger"`
	Info    struct {
		Description string `json:"info"`
		Title       string `json:"Title"`
		Version     string `json:"Version"`
	} `json:"info"`
	Host     string `json:"host"`
	BasePath string `json:"basePath"`
	Paths    map[string]map[string]struct {
		Security    []map[string]interface{} `json:"security"`
		Summary     string                   `json:"summary"`
		Description string                   `json:"description"`
		Parameters  []parameters             `json:"parameters"`
	} `json:"paths"`
	Definitions map[string]struct {
		Properties map[string]struct {
			Type string `json:"type"`
			Ref  string `json:"$ref"`
		} `json:"properties"`
	} `json:"definitions"`
	SecurityDefinitions map[string]struct {
		Type string `json:"Type"`
		Name string `json:"Name"`
		In   string `json:"In"`
	} `json:"securityDefinitions"`
}

type Swagger struct {
	inpt dtos.InputSwagger
}

func (c *Swagger) Read(args ...string) error {
	if err := validate.InputErrors(args, &c.inpt); err != nil {
		return err
	}

	return nil
}

func (c *Swagger) Eval() error {
	data := SwaggerData{}

	result, err := os.ReadFile(c.inpt.Path)
	if err != nil {
		return errors.BadRequest("Something wrong with the path provided")
	}

	if err := json.Unmarshal(result, &data); err != nil {
		fmt.Println(err)
	}

	repowork, err := repositorywork.NewWorkspaceRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to workspace table!")
	}

	repoenvs, err := repositoryenvs.NewEnvsRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to envs table!")
	}

	workspace := repowork.FindByName(c.inpt.Parent)

	resources := []resourcesdtos.InputUpdate{}

	for endpoint, item := range data.Paths {
		for method, v := range item {

			headers := c.GenHeaders(v.Parameters)

			for _, sec := range v.Security {
				for k := range sec {
					value := strings.ToUpper(fmt.Sprintf("%s_%s", workspace.Name, k))

					if _, err := repoenvs.FindByKey(value); err == nil {
						value = "$" + value
						headers[k] = value
					}
				}

			}
			resource := resourcesdtos.InputUpdate{
				Name:     fmt.Sprintf("%s %s", strings.ToUpper(method), endpoint),
				Endpoint: endpoint,
				Method:   strings.ToUpper(method),
				Params:   c.GenParams(v.Parameters),
				Headers:  headers,
				Body:     nil,
			}

			resources = append(resources, resource)
		}

	}

	value := workspacedtos.InputUpdate{
		Resources: resources,
	}

	repowork.Update(workspace, &value)

	return nil
}

func (c *Swagger) GenParams(params []parameters) []resourcesdtos.KeyValue {
	data := []resourcesdtos.KeyValue{}

	if len(params) == 0 {
		return data
	}

	for _, item := range params {
		param := resourcesdtos.KeyValue{}
		if "query" == item.In {
			param[item.Name] = ""

			data = append(data, param)
		}
	}

	return data
}

func (c *Swagger) GenHeaders(params []parameters) resourcesdtos.KeyValue {
	data := resourcesdtos.KeyValue{}

	if len(params) == 0 {
		return data
	}

	for _, item := range params {
		if "header" == item.In {
			data[item.Name] = ""
		}
	}

	return data
}

func (c *Swagger) Print() error {
	msg := fmt.Sprintf("[%s] command created with success!\n", "c.cmd.Name")

	fmt.Println(msg)

	return nil
}

func (w *Swagger) Run(args ...string) error {
	if err := w.Read(args...); err != nil {
		return err
	}

	if err := w.Eval(); err != nil {
		return err
	}

	if err := w.Print(); err != nil {
		return err
	}

	return nil
}

func SwaggerInit() repl.Repl {
	return &Swagger{}
}

var exemple = `

`
