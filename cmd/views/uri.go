package views

import (
	"github.com/jroimartin/gocui"
)

func Uri(g *gocui.Gui) error {
	maxX, _ := g.Size()

	if v, err := g.SetView("uri", 4, 2, maxX-4, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Uri "
	}

	return nil

}
