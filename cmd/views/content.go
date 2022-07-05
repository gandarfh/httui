package views

import "github.com/jroimartin/gocui"

func Content(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("content", 32, 5, maxX-4, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = " Content "

	}

	return nil
}
