package views

import (
	"fmt"

	"github.com/gandarfh/httui/cmd"
	"github.com/jroimartin/gocui"
)

func Tabs(g *gocui.Gui, config *cmd.Config) error {
	maxX, _ := g.Size()

	x, x1 := 38, (maxX/2)+10
	y, y1 := 5, 7

	if _, err := g.SetView("tabs", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := tabBindings(g); err != nil {
			return err
		}
		Tab("body", x+2, x+8, g, config)
		Tab("params", x+9, x+17, g, config)
		Tab("header", x+18, x+26, g, config)
	}

	return nil
}

func Tab(name string, x int, x1 int, g *gocui.Gui, config *cmd.Config) error {
	y, y1 := 5, 7

	if v, err := g.SetView("tab-"+name, x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v, name)

		v.Frame = false
	}

	return nil
}

var tabsViews = []string{"tab-body", "tab-params", "tab-header"}
var tabsActive = 0

func tabBindings(g *gocui.Gui) error {
	for _, item := range append(tabsViews, "tabs") {

		if err := g.SetKeybinding(item, gocui.KeyTab, gocui.ModNone, move); err != nil {

			return err
		}

	}

	return nil
}

func move(g *gocui.Gui, v *gocui.View) error {
	cmd.Navigate("next", tabsViews, &active)(g, v)

	v.Highlight = true
	v.SelFgColor = gocui.ColorGreen

	return nil
}
