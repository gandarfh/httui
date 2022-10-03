package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gandarfh/maid-san/internal/commands/resources/dtos"
	"github.com/gandarfh/maid-san/internal/commands/resources/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/validate"
	"github.com/google/uuid"
)

var (
	tmp_file = filepath.Join(os.TempDir(), uuid.New().String()+".json")
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

	if err := c.create_tmp_file(resource); err != nil {
		return err
	}

	cmd := exec.Command("lvim", tmp_file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()

	file, err := os.ReadFile(tmp_file)
	if err != nil {
		return errors.BadRequest("Error when try read the file.\n", err.Error())
	}

	value := dtos.InputUpdate{}
	if err := json.Unmarshal(file, &value); err != nil {
		return errors.BadRequest("Error when try marshal to repository.\n", err.Error())
	}

	repo.Update(resource, &value)

	if err := os.Remove(tmp_file); err != nil {
		return errors.BadRequest("Error when try delete the tmp file.\n", err.Error())
	}

	return nil
}

func (c *Update) create_tmp_file(resource *repository.Resources) error {
	params := []dtos.KeyValue{}
	for _, param := range resource.Params {
		params = append(params, dtos.KeyValue{Key: param.Key, Value: param.Value})
	}

	headers := []dtos.KeyValue{}
	for _, header := range resource.Headers {
		headers = append(headers, dtos.KeyValue{Key: header.Key, Value: header.Value})
	}

	data := dtos.InputUpdate{
		WorkspacesId: resource.WorkspacesId,
		Name:         resource.Name,
		Endpoint:     resource.Endpoint,
		Method:       resource.Method,
		Params:       params,
		Headers:      headers,
		Body:         resource.Body,
	}

	file, err := os.Create(tmp_file)
	if err != nil {
		return errors.BadRequest("Error when try create the file.\n", err.Error())
	}

	text, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return errors.BadRequest("Error when try marshal resource data to json file.\n", err.Error())
	}

	file.Write([]byte(text))

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
