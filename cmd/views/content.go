package views

import (
	"github.com/gandarfh/httui/cmd"
	"github.com/jroimartin/gocui"
)

func Content(g *gocui.Gui, config *cmd.Config) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("content", 32, 5, maxX-4, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = " Content "

	}

	return nil
}
