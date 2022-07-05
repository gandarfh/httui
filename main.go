package main

import (
	"log"

	"github.com/gandarfh/httui/cmd"
	v "github.com/gandarfh/httui/cmd/views"
	"github.com/jroimartin/gocui"
)

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.InputEsc = true

	g.SetManagerFunc(v.Layouts)

	if err := cmd.Bindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
