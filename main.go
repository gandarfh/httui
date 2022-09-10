package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gandarfh/httui-repl/internal/methods"
	"github.com/gandarfh/httui-repl/internal/methods/errors"
	"github.com/gandarfh/httui-repl/pkg/process"
)

const (
	version = "0.1.0"
)

func main() {
	wellcome()
	console()
}

func console() {
	tokens, err := start()
	if err != nil {
		log.Fatal(err)
	}

	err = process.Start(tokens, methods.Commands)
	if err != nil {
		command := errors.Init(err)

		command.Read(tokens...)
		command.Eval()
		command.Print()
	}

	console()
}

func start() ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("httui=> ")
	cmd, err := reader.ReadString('\n')

	if err != nil {
		return nil, fmt.Errorf("Error when try marshal input data to string!")
	}

	tokens := strings.Split(strings.TrimSpace(cmd), " ")

	return tokens, nil
}

func wellcome() {
	fmt.Printf("Wellcome my love >.<  ───  v%s.\n", version)
	fmt.Println("Want more?")
	fmt.Print("\n")
}
