package main

import (
	"log"

	"github.com/gandarfh/httui/cmd"
	v "github.com/gandarfh/httui/cmd/views"
	"github.com/jroimartin/gocui"
)

func main() {
	config, err := cmd.Connect()
	Error(err)

	g, err := gocui.NewGui(gocui.OutputNormal)
	Error(err)

	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.InputEsc = true

	g.SetManagerFunc(Layouts(config))

	err = cmd.Bindings(g)
	Error(err)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func Layouts(config *cmd.Config) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {

		if err := v.Uri(g, config); err != nil {
			return err
		}

		if err := v.Endpoints(g, config); err != nil {
			return err
		}

		if err := v.Content(g, config); err != nil {
			return err
		}

		return nil

	}
}

func Error(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
