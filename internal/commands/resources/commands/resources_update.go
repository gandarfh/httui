package commands

import (
	"fmt"

	"github.com/gandarfh/maid-san/internal/commands/resources/dtos"
	"github.com/gandarfh/maid-san/internal/commands/resources/repository"
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
	repo, err := repository.NewResourcesRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	resource := repo.FindByName(c.inpt.Name)

	preview := vim.NewPreview(resource)
	defer preview.Close()

	if err := preview.Open(); err != nil {
		return err
	}

	value := dtos.InputUpdate{}
	if err := preview.Execute(&value); err != nil {
		return err
	}

	repo.Update(resource, &value)

	return nil
}

func (c *Update) Print() error {
	msg := fmt.Sprintf("[%s] Resource success updated!\n", c.inpt.Name)

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
	repo, _ := repository.NewResourcesRepo()
	list := repo.List()

	commands := repl.CommandList{}

	commands = append(commands, repl.Command{
		Key:  "update",
		Repl: UpdateInit(),
	})

	for _, item := range *list {
		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf("%s %s", "update", item.Parent().Name),
			Repl: UpdateInit(),
		})

		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf(`%s %s name="%s"`, "update", item.Parent().Name, item.Name),
			Repl: UpdateInit(),
		})
	}

	return commands
}

func UpdateInit() repl.Repl {
	return &Update{}
}
