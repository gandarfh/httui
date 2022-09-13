<div align="center">
  <img
    src="assets/logo.png"
    alt="logo"
    style="width: 360px; height: 100%"
  />
  <h1 style="margin-top: 16px;">Maid-san</h1>
  <p>My CLI with all my necessary tools for work.</p>
</div>

## About the Project

This project will center all my needed tools like: http client, db view, healthcheck micro-services, etc.

### Tech Stack

- [Golang](https://go.dev/)
- [Sqlite](https://www.sqlite.org/index.html)

### Features

- Http client.
  - Multi workspaces with infinites resources.
  - Variables to use in values.
  - Open response into vim editor.
- Connect with dbui to connect with sql databases.
- Connect with mngr to connect with mongo databases.
- Calendar to connect with my google calendar and Alexa
- Review pull requests from github.

## Contact

Twitter - [@gandarfh](https://twitter.com/gandarfh)

<!-- Acknowledgments -->

## Acknowledgements

- [Validator](github.com/go-playground/validator/v10)
- [Liner](github.com/peterh/liner)
- [Gorm](gorm.io/gorm)

## Guide

When create a new command, follow this simple structure:

```go
package commandname

import (
	"github.com/gandarfh/maid-san/pkg/repl"
)

type CommandName struct{}

func (w *CommandName) Read(args ...string) error {
	return nil
}

func (w *CommandName) Eval() error {
	return nil
}

func (w *CommandName) Print() error {
	return nil
}

func (w *CommandName) Run(args ...string) error {
	w.Read(args...)
	w.Eval()
	w.Print()

	return nil
}

func Init() repl.Repl {
	return &CommandName{}
}
```
