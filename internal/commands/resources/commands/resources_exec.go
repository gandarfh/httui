package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gandarfh/maid-san/internal/commands/resources/repository"
	"github.com/gandarfh/maid-san/pkg/client"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/hokaccha/go-prettyjson"
	"github.com/logrusorgru/aurora/v3"
)

type Exec struct {
	ResourceId uint
	resource   *repository.Resources
	data       map[string]any
}

func (c *Exec) Read(args ...string) error {
	args = strings.Split(args[0], " ")
	resourceId, err := strconv.Atoi(args[2])

	if err != nil {
		return errors.UnprocessableEntity("Id provided isn't uint")
	}

	c.ResourceId = uint(resourceId)

	return nil
}

type teste struct{}

func (c *Exec) Eval() error {
	repo, err := repository.NewResourcesRepo()

	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	c.resource = repo.Find(c.ResourceId)
	workspace := c.resource.Parent()

	url := workspace.Uri + c.resource.Endpoint
	res := client.Request(url, c.resource.Method).Body(c.resource.Body)

	for _, item := range c.resource.Headers {
		res.Header(item.Key, item.Value)
	}

	for _, item := range c.resource.Params {
		res.Params(item.Key, item.Value)
	}

	data, err := res.Decode(&c.data)

	if err != nil {
		return errors.BadRequest()
	}

	fmt.Printf("%s - %s - %s \n", aurora.Yellow(c.resource.Method).Bold(), aurora.Bold(url), aurora.Green(data.Status).Bold())

	return nil
}

func (c *Exec) Print() error {
	for _, item := range c.resource.Headers {
		fmt.Println(aurora.Cyan(item.Key).String()+":", item.Value)
	}
	fmt.Print("\n")

	response, _ := prettyjson.Marshal(c.data)

	fmt.Println(string(response))
	fmt.Print("\n")

	return nil
}

func (w *Exec) Run(args ...string) error {
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

func ExecSubs() repl.CommandList {
	repo, _ := repository.NewResourcesRepo()
	list := repo.List()

	commands := repl.CommandList{}

	for _, item := range *list {
		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf("%s %d", "exec", item.ID),
			Repl: ExecInit(),
		})

	}

	return commands
}

func ExecInit() repl.Repl {
	return &Exec{}
}
