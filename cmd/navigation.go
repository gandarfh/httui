package cmd

import "github.com/jroimartin/gocui"

var (
	viewArr = []string{"endpoints", "content", "uri"}
	active  = 0
)

func SetActive(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func Navigate(direction string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		var index int
		var name string

		if direction == "next" {
			index = (active + 1) % len(viewArr)
			name = viewArr[index]
		}

		if direction == "prev" {
			index = (active + len(viewArr) - 1) % len(viewArr)
			name = viewArr[index]
		}

		if _, err := g.View("endpoints"); err != nil {
			return err
		}

		if _, err := SetActive(g, name); err != nil {
			return err
		}

		active = index
		return nil
	}

}
