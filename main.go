package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/peterh/liner"

	"github.com/gandarfh/httui-repl/internal/commands"
	"github.com/gandarfh/httui-repl/internal/commands/errors"
	"github.com/gandarfh/httui-repl/pkg/process"
	"github.com/gandarfh/httui-repl/pkg/utils"
)

const (
	version    = "0.1.0"
	history_fn = "./.httui_history"
)

func main() {
	line := liner.NewLiner()
	defer line.Close()

	fmt.Printf("Wellcome my love >.<  ───  v%s.\n", version)
	fmt.Println("Want more?")
	fmt.Print("\n")

	console(line)

}

func console(line *liner.State) {
	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	if output, err := line.Prompt("httui=> "); err == nil {
		line.AppendHistory(output)
		args := utils.SplitArgs(strings.TrimSpace(output))

		err = process.Start(args, commands.Commands)
		if err != nil {
			command := errors.Init(err).(*errors.Error)
			command.Run(args...)
		}

	} else {
		log.Print("Error reading line: ", err)
		os.Exit(0)
	}

	if f, err := os.Create(history_fn); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}

	console(line)
}
