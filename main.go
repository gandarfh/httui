package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/peterh/liner"

	"github.com/gandarfh/httui-repl/internal/methods"
	"github.com/gandarfh/httui-repl/internal/methods/errors"
	"github.com/gandarfh/httui-repl/pkg/process"
)

const (
	version    = "0.1.0"
	history_fn = "./.httui_history"
)

func main() {
	line := liner.NewLiner()
	defer line.Close()

	wellcome()
	console(line)

}

func console(line *liner.State) {
	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	if value, err := line.Prompt("httui=> "); err == nil {
		line.AppendHistory(value)
		args, _ := getArgs(value)

		err = process.Start(args, methods.Commands)
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

func getArgs(command string) ([]string, error) {
	tokens := strings.Split(strings.TrimSpace(command), " ")

	return tokens, nil
}

func wellcome() {
	fmt.Printf("Wellcome my love >.<  ───  v%s.\n", version)
	fmt.Println("Want more?")
	fmt.Print("\n")
}
