package main

import (
	"log"

	"github.com/gandarfh/httui/cmd"
	v "github.com/gandarfh/httui/cmd/views"
	"github.com/jroimartin/gocui"
)

func layouts(g *gocui.Gui) error {
	if err := v.Endpoints(g); err != nil {
		return err
	}

	if err := v.Uri(g); err != nil {
		return err
	}

	if err := v.Content(g); err != nil {
		return err
	}

	return nil
}

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layouts)

	if err := cmd.Bindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
