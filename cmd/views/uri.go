package views

import (
	"fmt"

	"github.com/gandarfh/httui/cmd"
	"github.com/gandarfh/httui/cmd/model"
	"github.com/jroimartin/gocui"
	c "github.com/logrusorgru/aurora/v3"
)

func getCurrentUri(v *gocui.View, name string) {
	fmt.Fprintln(v, c.Bold(c.Yellow(name)))
}

func Uri(g *gocui.Gui) error {
	maxX, _ := g.Size()

	if v, err := g.SetView("uri", 4, 2, maxX-4, 4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if _, err = cmd.SetActive(g, "uri"); err != nil {
			return err
		}

		if err := uBindings(g); err != nil {
			return err
		}

		Bus.Subscribe("uri:current-uri", getCurrentUri)
		v.Title = " Uri "

	}

	return nil
}

func NewUri(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("")

	maxX, maxY := g.Size()

	if v, err := g.SetView("create-uri", maxX/5, (maxY / 2), (maxX - (maxX / 5)), (maxY/2)+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = " Create uri "
		v.Editable = true
		cmd.SetActive(g, "create-uri")
	}

	return nil
}

func ListUris(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("")
	maxX, maxY := g.Size()

	if v, err := g.SetView("select-uri", maxX/5, (maxY / 3), (maxX - (maxX / 5)), (maxY/2)+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if list, err := model.GetUris(); err != nil {
			return err
		} else {
			for _, name := range list {
				fmt.Fprintln(v, "", c.Bold(c.White(name)))
			}
		}

		v.Title = " Select uri "
		cmd.SetActive(g, "select-uri")
	}

	return nil
}

func uBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("uri", 'c', gocui.ModNone, NewUri); err != nil {
		return err
	}

	if err := g.SetKeybinding("uri", 's', gocui.ModNone, ListUris); err != nil {
		return err
	}

	if err := g.SetKeybinding("create-uri", gocui.KeyEnter, gocui.ModNone, saveNewUri); err != nil {
		return err
	}

	if err := g.SetKeybinding("create-uri", gocui.KeyCtrlC, gocui.ModNone, cmd.Close("create-uri", "uri")); err != nil {
		return err
	}

	if err := g.SetKeybinding("select-uri", 'j', gocui.ModNone, cmd.CursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("select-uri", 'k', gocui.ModNone, cmd.CursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("select-uri", gocui.KeyCtrlC, gocui.ModNone, cmd.Close("select-uri", "uri")); err != nil {
		return err
	}
	if err := g.SetKeybinding("select-uri", gocui.KeyEnter, gocui.ModNone, selectUri); err != nil {
		return err
	}

	return nil
}

func selectUri(g *gocui.Gui, v *gocui.View) error {
	defer cmd.Close("select-uri", "uri")(g, v)

	uView, err := g.View("uri")
	if err != nil {
		return err
	}

	line, err := cmd.GetLine(v)

	if err != nil {
		return err
	}

	eView, err := g.View("endpoints")

	if err != nil {
		return err
	}
	eView.Clear()
	uView.Clear()
	Bus.Publish("uri:current-uri", uView, line)
	Bus.Publish("endpoints:get", g, eView, line)
	return nil
}

func saveNewUri(g *gocui.Gui, v *gocui.View) error {
	defer cmd.Close("create-uri", "uri")(g, v)
	name, err := v.Line(0)

	if err != nil {
		return err
	}

	model.CreateUri(name)

	return nil
}
