package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gandarfh/maid-san/internal/commands/generator/dtos"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/utils"
	"github.com/gandarfh/maid-san/pkg/validate"
)

var (
	templatepath = "./internal/commands/command/commands/template"
)

type Create struct {
	cmd    dtos.InputCreate
	fields map[string]any
}

func (c *Create) Read(args ...string) error {
	if err := validate.InputErrors(args, &c.cmd); err != nil {
		return err
	}

	fields, err := utils.ArgsFormat(args[1:])

	if err != nil {
		return errors.InternalServer()
	}

	c.fields = fields

	return nil
}

func (c *Create) Eval() error {
	if err := c.CreateFileByTemplate(c.cmd.Template, c.cmd.Path); err != nil {
		return err
	}

	return nil
}

func (c *Create) Print() error {
	msg := fmt.Sprintf("[%s] command created with success!\n", "c.cmd.Name")

	fmt.Println(msg)

	return nil
}

func (w *Create) Run(args ...string) error {
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

func (c *Create) CreateFileByTemplate(template, path string) error {
	files, err := os.ReadDir(template)
	if err != nil {
		return errors.BadRequest("Error when try read template dir.\n", err.Error())
	}

	for _, item := range files {
		subdirtemplate := filepath.Join(template, item.Name())
		subpath := filepath.Join(path, item.Name())

		if item.IsDir() {
			if err := os.MkdirAll(subpath, 0777); err != nil {
				return errors.BadRequest("Error when try create the directory.\n", err.Error())
			}

			c.CreateFileByTemplate(subdirtemplate, subpath)
			continue
		}

		filename := fmt.Sprintf("%s%s", subpath, c.cmd.Type)

		for key, value := range c.fields {
			old := fmt.Sprintf("{{%s}}", key)
			filename = strings.ReplaceAll(filename, old, string(value.(string)))
		}

		file, err := os.Create(filename)
		if err != nil {
			return errors.BadRequest("Error when try create the file.\n", err.Error())
		}

		data, err := os.ReadFile(subdirtemplate)
		if err != nil {
			return errors.BadRequest("Error when try read template.\n", err.Error())
		}

		text := string(data)

		for key, value := range c.fields {
			old := fmt.Sprintf("{{%s}}", key)
			text = strings.ReplaceAll(text, old, string(value.(string)))
		}

		file.Write([]byte(text))
	}

	return nil
}

func CreateInit() repl.Repl {
	return &Create{}
}
