package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gandarfh/maid-san/internal/commands/resources/dtos"
	"github.com/gandarfh/maid-san/internal/commands/resources/repository"
	"github.com/gandarfh/maid-san/pkg/client"
	"github.com/gandarfh/maid-san/pkg/errors"
	"github.com/gandarfh/maid-san/pkg/repl"
	"github.com/gandarfh/maid-san/pkg/truncate"
	"github.com/gandarfh/maid-san/pkg/utils"
	"github.com/gandarfh/maid-san/pkg/validate"
	"github.com/gandarfh/maid-san/pkg/vim"
	"github.com/hokaccha/go-prettyjson"
	"github.com/logrusorgru/aurora/v3"
)

type Exec struct {
	inpt     dtos.InputExec
	resource *repository.Resources
	data     map[string]any
	withVim  bool
}

func (c *Exec) Read(args ...string) error {
	if err := validate.InputErrors(args, &c.inpt); err != nil {
		return err
	}

	c.withVim = strings.Contains(args[0], "vim")

	return nil
}

func (c *Exec) Eval() error {
	repo, err := repository.NewResourcesRepo()

	if err != nil {
		return errors.InternalServer("Error when connect to database!")
	}

	c.resource = repo.FindByName(c.inpt.Name)
	workspace := c.resource.Parent()

	url := utils.ReplaceByEnv(workspace.Uri) + utils.ReplaceByEnv(c.resource.Endpoint)

	fmt.Println()

	res := client.Request(url, c.resource.Method)

	rawbody, _ := c.resource.Body.MarshalJSON()
	bodystring := utils.ReplaceByEnv(string(rawbody))

	body := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodystring), &body); err != nil {
		panic(err)
	}

	if len(body) == 0 {
		res.Body(nil)
	} else {
		res.Body([]byte(bodystring))
	}

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
		if err := c.preview(c.resource, url, data.Status); err != nil {
			return err
		}
	}

	return nil
}

func (c *Exec) Print() error {
	for _, item := range c.resource.Headers {
		fmt.Println(aurora.Cyan(item.Key).String()+":", truncate.Dots(utils.ReplaceByEnv(item.Value), 20))
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

func (c *Exec) preview(resource *repository.Resources, url string, status string) error {
	params := dtos.KeyValue{}
	for _, item := range resource.Params {
		params[item.Key] = item.Value
	}

	headers := dtos.KeyValue{}
	for _, item := range resource.Headers {
		headers[item.Key] = item.Value
	}

	data := struct {
		Url     string        `json:"url"`
		Method  string        `json:"method"`
		Status  string        `json:"status"`
		Params  dtos.KeyValue `json:"params"`
		Headers dtos.KeyValue `json:"headers"`
		Body    any           `json:"body"`
	}{
		Url:     url,
		Method:  resource.Method,
		Status:  status,
		Params:  params,
		Headers: headers,
		Body:    c.data,
	}

	preview := vim.NewPreview(data)
	defer preview.Close()

	if err := preview.Open(); err != nil {
		return err
	}

	return nil
}

func ExecSubs() repl.CommandList {
	repo, _ := repository.NewResourcesRepo()
	list := repo.List()

	commands := repl.CommandList{}

	commands = append(commands, repl.Command{
		Key:  "exec",
		Repl: ExecInit(),
	})

	commands = append(commands, repl.Command{
		Key:  "vim exec",
		Repl: ExecInit(),
	})

	for _, item := range *list {
		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf(`%s %s`, "exec", item.Parent().Name),
			Repl: ExecInit(),
		})

		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf(`%s %s`, "vim exec", item.Parent().Name),
			Repl: ExecInit(),
		})

		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf(`%s %s name="%s"`, "exec", item.Parent().Name, item.Name),
			Repl: ExecInit(),
		})

		commands = append(commands, repl.Command{
			Key:  fmt.Sprintf(`%s %s name="%s"`, "vim exec", item.Parent().Name, item.Name),
			Repl: ExecInit(),
		})
	}

	return commands
}

func ExecInit() repl.Repl {
	return &Exec{}
}
