package views

import (
	"github.com/gandarfh/httui/cmd"
	"github.com/jroimartin/gocui"
)

func Body(g *gocui.Gui, config *cmd.Config) error {
	maxX, maxY := g.Size()

	x, x1 := 36, (maxX/3)+8
	y, y1 := 5, maxY-4

	if v, err := g.SetView("body", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := bBindings(g, config); err != nil {
			return err
		}

		v.Title = " Body "

		cmd.Bus.Subscribe("test", teste(v))
	}

	return nil
}

func bBindings(g *gocui.Gui, config *cmd.Config) error {

	if err := g.SetKeybinding("body", 'i', gocui.ModNone, insertMode); err != nil {
		return err
	}

	if err := g.SetKeybinding("body", gocui.KeyEsc, gocui.ModNone, normalMode); err != nil {
		return err
	}

	return nil
}

func insertMode(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("")
	g.DeleteKeybinding("body", 'i', gocui.ModNone)
	v.Editable = true

	return nil
}

func normalMode(g *gocui.Gui, v *gocui.View) error {
	v.Editable = false
	if err := g.SetKeybinding("body", 'i', gocui.ModNone, insertMode); err != nil {
		return err
	}
	if err := cmd.Bindings(g); err != nil {
		return err
	}

	return nil
}
