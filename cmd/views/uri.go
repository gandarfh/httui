package views

import (
	"fmt"
	"strings"

	"github.com/gandarfh/httui/cmd"
	"github.com/gandarfh/httui/cmd/model"
	"github.com/jroimartin/gocui"
	c "github.com/logrusorgru/aurora/v3"
)

func Uri(g *gocui.Gui, config *cmd.Config) error {
	maxX, _ := g.Size()

	if v, err := g.SetView("uri", 4, 2, maxX/3, 4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if _, err = cmd.SetActive(g, "uri"); err != nil {
			return err
		}

		if err := uBindings(g, config); err != nil {
			return err
		}

		uDefault := config.GetDefaultUri()

		fmt.Fprintln(v, c.Bold(c.Yellow(uDefault.Alias)))

		cmd.Bus.Subscribe("uri:current-uri", getCurrentUri)
		v.Title = " Uri "
	}

	return nil
}

func newUri(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("")

	maxX, maxY := g.Size()

	y, y1 := (maxY / 3), (maxY/3)+2
	x, x1 := maxX/3, (maxX - (maxX / 3))

	if v, err := g.SetView("create-uri", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = " Create uri "
		v.Editable = true
		cmd.SetActive(g, "create-uri")
	}

	return nil
}

func listUris(config *cmd.Config) func(g *gocui.Gui, v *gocui.View) error {

	return func(g *gocui.Gui, v *gocui.View) error {

		g.DeleteKeybindings("")
		maxX, maxY := g.Size()

		y, y1 := (maxY / 3), (maxY / 2)
		x, x1 := maxX/3, (maxX - (maxX / 3))

		if v, err := g.SetView("select-uri", x, y, x1, y1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			list := config.Uris

			for _, item := range *list {
				fmt.Fprintln(v, "", c.Bold(c.White(item.Alias)))
			}

			v.Title = " Select uri "
			cmd.SetActive(g, "select-uri")
		}

		return nil

	}
}

func uBindings(g *gocui.Gui, config *cmd.Config) error {
	if err := g.SetKeybinding("uri", 'c', gocui.ModNone, newUri); err != nil {
		return err
	}

	if err := g.SetKeybinding("uri", 's', gocui.ModNone, listUris(config)); err != nil {
		return err
	}

	if err := g.SetKeybinding("create-uri", gocui.KeyEnter, gocui.ModNone, saveNewUri(config)); err != nil {
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

func getCurrentUri(v *gocui.View, name string) {
	fmt.Fprintln(v, c.Bold(c.Yellow(name)))
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

	uri := strings.TrimLeft(line, " ")

	cmd.Bus.Publish("uri:current-uri", uView, uri)
	cmd.Bus.Publish("endpoints:get", g, eView, &uri)

	return nil
}

func saveNewUri(config *cmd.Config) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {

		defer cmd.Close("create-uri", "uri")(g, v)

		name, err := v.Line(0)

		if err != nil {
			return err
		}

		uri := model.Uri{Alias: name, Endpoints: &[]model.Endpoint{}}

		if err := config.CreateUri(&uri); err != nil {
			return err
		}

		return nil

	}
}
