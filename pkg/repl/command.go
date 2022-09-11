package repl

import "fmt"

type Command struct {
	Key         string
	Parent      string
	SubCommands SubCommands
	Repl        Repl
}

type SubCommands []Command
type CommandList []Command

func (c *SubCommands) Keys() []string {
	keys := []string{}

	for _, item := range *c {
		keys = append(keys, item.Key)
	}

	return keys
}

func (clist *CommandList) Find(key string) (Command, error) {
	finded := &Command{}

	for _, sub := range *clist {
		if sub.Key == key {
			finded = &sub
			break
		}
	}

	if finded == nil {
		return Command{}, fmt.Errorf("Not found!")
	}

	return *finded, nil
}
