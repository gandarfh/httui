package views

import (
	"github.com/gandarfh/httui/cmd"
	"github.com/jroimartin/gocui"
)

func Header(g *gocui.Gui, config *cmd.Config) error {
	maxX, maxY := g.Size()

	// if v, err := g.SetView("endpoints", 4, 5, 34, maxY-4); err != nil {
	// x, x1 := 4, 36
	// y, y1 := maxY/2+1, maxY-4

	x, x1 := 38, (maxX/2)+10
	y, y1 := 9, maxY/2

	if v, err := g.SetView("header", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		cmd.Bus.Subscribe("test", teste(v))
		v.Title = " Headers "
	}

	return nil
}
