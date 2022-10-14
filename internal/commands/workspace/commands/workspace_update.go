package commands

import (
	"fmt"

	resourcesdtos "github.com/gandarfh/maid-san/internal/commands/resources/dtos"
	"github.com/gandarfh/maid-san/internal/commands/workspace/dtos"
	"github.com/gandarfh/maid-san/internal/commands/workspace/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/validate"
	"github.com/gandarfh/maid-san/pkg/vim"
)

type Update struct {
	inpt dtos.InputUpdate
}

func (c *Update) Read(args ...string) error {
	if err := validate.InputErrors(args, &c.inpt); err != nil {
		return err
	}

	return nil
}

func (c *Update) Eval() error {
	repo, err := repository.NewWorkspaceRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	workspace := repo.FindByName(c.inpt.Name)

	resources := []resourcesdtos.InputUpdate{}
	for _, item := range workspace.Resources {

		params := []resourcesdtos.KeyValue{}
		for _, item := range item.Params {
			key := resourcesdtos.KeyValue{}
			key[item.Key] = item.Value
			params = append(params, key)
		}

		headers := resourcesdtos.KeyValue{}
		for _, item := range item.Headers {
			headers[item.Key] = item.Value
		}

		resources = append(resources, resourcesdtos.InputUpdate{
			WorkspacesId: item.WorkspacesId,
			Parent:       item.Parent().Name,
			Name:         item.Name,
			Endpoint:     item.Endpoint,
			Method:       item.Method,
			Params:       params,
			Headers:      headers,
			Body:         item.Body,
		})

	}

	data := dtos.InputUpdate{
		Name:      workspace.Name,
		Uri:       workspace.Uri,
		Resources: resources,
	}

	preview := vim.NewPreview(data)
	defer preview.Close()

	if err := preview.Open(); err != nil {
		return err
	}

	value := dtos.InputUpdate{}
	if err := preview.Execute(&value); err != nil {
		return err
	}

	repo.Update(workspace, &value)

	return nil
}

func (c *Update) Print() error {
	msg := fmt.Sprintf("[%s] Workspace success updated!\n", c.inpt.Name)

	fmt.Println(msg)

	return nil
}

func (w *Update) Run(args ...string) error {
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

func UpdateSubs() repl.CommandList {
	repo, _ := repository.NewWorkspaceRepo()
	list := repo.List()

	commands := repl.CommandList{}

	commands = append(commands, repl.Command{
		Key:  "update",
		Repl: UpdateInit(),
	})

	for _, item := range *list {
		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf(`%s name="%s"`, "update", item.Name),
			Repl: UpdateInit(),
		})
	}

	return commands
}

func UpdateInit() repl.Repl {
	return &Update{}
}
