package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"

	"github.com/gandarfh/maid-san/internal/commands"
	"github.com/gandarfh/maid-san/internal/commands/errors"
	"github.com/gandarfh/maid-san/internal/commands/welcome"
	"github.com/gandarfh/maid-san/pkg/process"
	"github.com/gandarfh/maid-san/pkg/utils"
)

const (
	history_fn = ".maid-san_history"
)

func main() {
	line := liner.NewLiner()
	defer line.Close()

	welcome := welcome.Init()
	welcome.Print()

	line.SetCompleter(func(l string) (c []string) {
		for _, i := range commands.Cmds() {
			if strings.HasPrefix(i.Key, strings.ToLower(l)) {
				c = append(c, i.Key)
			}
		}
		return
	})

	console(line)

}

func console(line *liner.State) {
	home, _ := os.UserHomeDir()
	if f, err := os.Open(filepath.Join(home, history_fn)); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	if output, err := line.Prompt("💕[maid-san]: "); err == nil {
		line.AppendHistory(output)
		args := utils.SplitArgs(strings.TrimSpace(output))

		err = process.Start(args, commands.Cmds())
		if err != nil {
			command := errors.Init(err)
			command.Run(args...)
		}

	} else {
		log.Print("Error reading line: ", err)
		os.Exit(0)
	}

	if f, err := os.Create(filepath.Join(home, history_fn)); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}

	console(line)
}
