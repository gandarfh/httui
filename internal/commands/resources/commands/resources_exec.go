package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gandarfh/maid-san/internal/commands/resources/dtos"
	"github.com/gandarfh/maid-san/internal/commands/resources/repository"
	"github.com/gandarfh/maid-san/pkg/client"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/utils"
	"github.com/hokaccha/go-prettyjson"
	"github.com/logrusorgru/aurora/v3"
)

type Exec struct {
	ResourceId uint
	resource   *repository.Resources
	data       map[string]any
	withVim    bool
}

func (c *Exec) Read(args ...string) error {
	var (
		resourceId = 0
		err        error
	)

	c.withVim = strings.Contains(args[0], "vim")
	args = strings.Split(args[0], " ")

	if c.withVim {
		resourceId, err = strconv.Atoi(args[3])
	} else {
		resourceId, err = strconv.Atoi(args[2])
	}

	if err != nil {
		return errors.UnprocessableEntity("Id provided isn't uint")
	}

	c.ResourceId = uint(resourceId)

	return nil
}

func (c *Exec) Eval() error {
	repo, err := repository.NewResourcesRepo()

	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	c.resource = repo.Find(c.ResourceId)
	workspace := c.resource.Parent()

	url := utils.ReplaceByEnv(workspace.Uri) + utils.ReplaceByEnv(c.resource.Endpoint)

	rawbody, _ := c.resource.Body.MarshalJSON()
	body := utils.ReplaceByEnv(string(rawbody))

	res := client.Request(url, utils.ReplaceByEnv(c.resource.Method)).Body(body)

	for _, item := range c.resource.Headers {
		res.Header(item.Key, utils.ReplaceByEnv(item.Value))
	}

	for _, item := range c.resource.Params {
		res.Params(item.Key, utils.ReplaceByEnv(item.Value))
	}

	data, err := res.Decode(&c.data)

	if err != nil {
		return errors.BadRequest()
	}

	fmt.Printf("%s - %s - %s \n", aurora.Yellow(c.resource.Method).Bold(), aurora.Bold(url), aurora.Green(data.Status).Bold())

	if c.withVim {
		if err := c.create_tmp_file(c.resource, url, data.Status); err != nil {
			return err
		}

		cmd := exec.Command("lvim", tmp_file)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()

		if err := os.Remove(tmp_file); err != nil {
			return errors.BadRequest("Error when try delete the tmp file.\n", err.Error())
		}
	}

	return nil
}

func (c *Exec) Print() error {
	for _, item := range c.resource.Headers {
		fmt.Println(aurora.Cyan(item.Key).String()+":", utils.ReplaceByEnv(item.Value))
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

func (c *Exec) create_tmp_file(resource *repository.Resources, url string, status string) error {
	file, err := os.Create(tmp_file)
	if err != nil {
		return errors.BadRequest("Error when try create the file.\n", err.Error())
	}

	params := []dtos.KeyValue{}
	for _, param := range resource.Params {
		params = append(params, dtos.KeyValue{Key: param.Key, Value: utils.ReplaceByEnv(param.Value)})
	}

	headers := []dtos.KeyValue{}
	for _, header := range resource.Headers {
		headers = append(headers, dtos.KeyValue{Key: header.Key, Value: utils.ReplaceByEnv(header.Value)})
	}

	data := struct {
		Url     string          `json:"url"`
		Method  string          `json:"method"`
		Status  string          `json:"status"`
		Params  []dtos.KeyValue `json:"params"`
		Headers []dtos.KeyValue `json:"headers"`
		Body    any             `json:"body"`
	}{
		Url:     url,
		Method:  resource.Method,
		Status:  status,
		Params:  params,
		Headers: headers,
		Body:    c.data,
	}

	text, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return errors.BadRequest("Error when try marshal resource data to json file.\n", err.Error())
	}

	file.Write([]byte(text))

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

		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf("%s %d", "vim exec", item.ID),
			Repl: ExecInit(),
		})
	}

	return commands
}

func ExecInit() repl.Repl {
	return &Exec{}
}
