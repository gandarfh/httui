package views

import (
	"github.com/gandarfh/httui/cmd"
	"github.com/jroimartin/gocui"
)

func Preview(g *gocui.Gui, config *cmd.Config) error {
	maxX, _ := g.Size()

	x, x1 := (maxX/3)+2, maxX-6
	y, y1 := 2, 4

	if v, err := g.SetView("preview", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = " Preview "

		cmd.Bus.Subscribe("test", teste(v))
	}

	return nil
}
