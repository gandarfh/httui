package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	notfound "github.com/gandarfh/httui-repl/cmd/methods/not_found"
	"github.com/gandarfh/httui-repl/cmd/process"
)

const (
	version = "0.1.0"
)

func main() {
	tokens, err := start()
	if err != nil {
		log.Fatal(err)
	}

	process, err := process.Start(tokens)
	if err != nil {
		process := notfound.Init()

		process.Read(tokens)
		process.Eval()
		process.Print()
	} else {
		// TODO: tratar erros para cada processo
		process.Read(tokens)
		process.Eval()
		process.Print()
	}
	main()
}

func wellcome() {
	fmt.Printf("Wellcome my love >.<  ───  v%s.\n", version)
	fmt.Println("Want more?")
	fmt.Print("\n")
}

func start() ([]string, error) {
	wellcome()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("httui=> ")
	cmd, err := reader.ReadString('\n')

	if err != nil {
		return nil, fmt.Errorf("Error when try marshal input data to string!")
	}

	tokens := strings.Split(strings.TrimSpace(cmd), " ")

	return tokens, nil
}
