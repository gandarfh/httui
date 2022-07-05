package cmd

import "github.com/jroimartin/gocui"

func SetActive(g *gocui.Gui, name string) (*gocui.View, error) {

	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}

	return g.SetViewOnTop(name)
}

func Navigate(direction string, views []string, active *int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		var index int
		var name string

		if direction == "next" {
			index = (*active + 1) % len(views)
			name = views[index]
		}

		if direction == "prev" {
			index = (*active + len(views) - 1) % len(views)
			name = views[index]
		}

		if _, err := SetActive(g, name); err != nil {
			return err
		}

		*active = index
		return nil
	}

}

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	maxLines := len(v.ViewBufferLines()) - 2

	if maxLines == cy {
		return nil
	}

	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// `delete` is what view you want remove from screen, and `active` what screen you want active.
// Create to use with bindings
// Close("create-uri", "uri")
func Close(delete string, active string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if err := g.DeleteView(delete); err != nil {
			return err
		}

		if err := Bindings(g); err != nil {
			return err
		}

		if _, err := SetActive(g, active); err != nil {
			return err
		}

		return nil
	}
}

func CloseList(list []string, active string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		for _, name := range list {
			if err := g.DeleteView(name); err != nil {
				return err
			}
		}

		if err := Bindings(g); err != nil {
			return err
		}

		if _, err := SetActive(g, active); err != nil {
			return err
		}

		return nil
	}
}

func GetLine(v *gocui.View) (string, error) {
	_, cy := v.Cursor()

	line, err := v.Line(cy)

	if err != nil {
		line = ""
		return line, err
	}

	return line, nil
}
