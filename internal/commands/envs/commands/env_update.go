package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gandarfh/maid-san/internal/commands/envs/dtos"
	"github.com/gandarfh/maid-san/internal/commands/envs/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/google/uuid"
)

var (
	tmp_file = filepath.Join(os.TempDir(), uuid.New().String()+".json")
)

type Update struct {
	inpt  dtos.InputUpdate
	EnvId uint
}

func (c *Update) Read(args ...string) error {
	args = strings.Split(args[0], " ")
	envId, err := strconv.Atoi(args[2])

	if err != nil {
		return errors.UnprocessableEntity("Id provided isn't uint")
	}

	c.EnvId = uint(envId)

	return nil
}

func (c *Update) Eval() error {
	repo, err := repository.NewEnvsRepo()
	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	env := repo.Find(c.EnvId)

	if err := c.create_tmp_file(env); err != nil {
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

	repo.Update(env, &value)

	if err := os.Remove(tmp_file); err != nil {
		return errors.BadRequest("Error when try delete the tmp file.\n", err.Error())
	}

	return nil
}

func (c *Update) create_tmp_file(env *repository.Envs) error {
	data := dtos.InputUpdate{
		Key:   env.Key,
		Value: env.Value,
	}
	c.inpt = data

	file, err := os.Create(tmp_file)
	if err != nil {
		return errors.BadRequest("Error when try create the file.\n", err.Error())
	}

	text, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return errors.BadRequest("Error when try marshal env data to json file.\n", err.Error())
	}

	file.Write([]byte(text))

	return nil
}

func (c *Update) Print() error {
	msg := fmt.Sprintf("[%s] Env success updated!\n", c.inpt.Key)

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
	repo, _ := repository.NewEnvsRepo()

	list := repo.List()

	commands := repl.CommandList{}

	for _, item := range *list {
		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf("%s %d", "update", item.ID),
			Repl: UpdateInit(),
		})
	}

	return commands
}

func UpdateInit() repl.Repl {

	return &Update{}
}