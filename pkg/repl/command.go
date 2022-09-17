package repl

import "fmt"

type Command struct {
	Key  string
	Repl Repl
}

type CommandList []Command

func (clist *CommandList) Find(key string) (Command, error) {
	finded := &Command{}

	for _, sub := range *clist {
		if sub.Key == key {
			finded = &sub
			break
		}
		finded = nil
	}

	if finded == nil {
		return Command{}, fmt.Errorf("Not found!")
	}

	return *finded, nil
}
