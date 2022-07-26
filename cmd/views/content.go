package views

import (
	"fmt"

	"github.com/gandarfh/httui/cmd"
	"github.com/jroimartin/gocui"
)

func Content(g *gocui.Gui, config *cmd.Config) error {
	if err := Tabs(g, config); err != nil {
		return err
	}
	if err := Params(g, config); err != nil {
		return err
	}
	if err := Body(g, config); err != nil {
		return err
	}
	if err := Header(g, config); err != nil {
		return err
	}
	if err := Response(g, config); err != nil {
		return err
	}
	if err := Preview(g, config); err != nil {
		return err
	}

	return nil
}

func teste(v *gocui.View) func(value string) {
	return func(value string) {

		fmt.Fprintln(v, value)
	}
}
