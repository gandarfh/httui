package elements

import (
	"github.com/jroimartin/gocui"
)

func Input(name string, x int, y int, x1 int, y1 int) func(g *gocui.Gui, parentView *gocui.View) error {
	return func(g *gocui.Gui, parentView *gocui.View) error {

		view, err := g.SetView(name, x, y, x1, y1)
		if err != gocui.ErrUnknownView {
			return err
		}

		view.Autoscroll = false
		view.Editable = true
		view.Title = name

		return nil
	}
}
