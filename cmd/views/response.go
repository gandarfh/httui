package views

import (
	"github.com/gandarfh/httui/cmd"
	"github.com/jroimartin/gocui"
)

func Response(g *gocui.Gui, config *cmd.Config) error {
	maxX, maxY := g.Size()

	x, x1 := (maxX/2)+12, maxX-6
	y, y1 := 5, maxY-4

	if v, err := g.SetView("response", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = " Response "

		cmd.Bus.Subscribe("test", teste(v))
	}

	return nil
}
