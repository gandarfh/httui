package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/generator/dtos"
	"github.com/gandarfh/maid-san/pkg/client"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/validate"
)

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
		Parameters  []struct {
			Name   string `json:"name"`
			In     string `json:"in"` // "body" | "query"
			Schema struct {
				Ref string `json:"$ref"`
			} `json:"schema"`
		} `json:"parameters"`
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
	teste := SwaggerData{}

	_, err := client.Get("http://localhost:5000/docs/doc.json").Decode(&teste)

	if err != nil {
		fmt.Println(err)
	}

	return nil
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
