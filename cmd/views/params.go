package views

import (
	"github.com/gandarfh/httui/cmd"
	"github.com/jroimartin/gocui"
)

func Params(g *gocui.Gui, config *cmd.Config) error {
	maxX, maxY := g.Size()

	// x, x1 := 36, (maxX/3)+8
	// y, y1 := maxY/2+1, maxY-4

	x, x1 := 38, (maxX/3)+8
	y, y1 := 9, maxY/2

	if v, err := g.SetView("params", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "                        Params "

	}

	return nil
}
