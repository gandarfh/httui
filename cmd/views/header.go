package views

import (
	"github.com/gandarfh/httui/cmd"
	"github.com/jroimartin/gocui"
)

func Header(g *gocui.Gui, config *cmd.Config) error {
	maxX, maxY := g.Size()

	x, x1 := (maxX/3)+10, (maxX/2)+18
	y, y1 := 5, maxY-4

	if v, err := g.SetView("header", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		cmd.Bus.Subscribe("test", teste(v))
		v.Title = " Headers "
	}

	return nil
}
