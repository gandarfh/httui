package cmd

import "github.com/jroimartin/gocui"

var views = []string{"uri", "endpoints", "content"}
var active = 0

func Bindings(g *gocui.Gui) error {
	// global key bindings
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'l', gocui.ModNone, Navigate("next", views, &active)); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'h', gocui.ModNone, Navigate("prev", views, &active)); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	if currentView := g.CurrentView().Name(); currentView == "create-uri" {
		return nil
	}

	return gocui.ErrQuit
}
