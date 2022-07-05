package views

import (
	"fmt"

	"github.com/gandarfh/httui/cmd"
	"github.com/gandarfh/httui/cmd/elements"
	"github.com/gandarfh/httui/cmd/model"
	"github.com/jroimartin/gocui"
	c "github.com/logrusorgru/aurora/v3"
)

func Endpoints(g *gocui.Gui) error {
	_, maxY := g.Size()

	if v, err := g.SetView("endpoints", 4, 5, 30, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := eBindings(g); err != nil {
			return err
		}

		v.Autoscroll = true
		g.Cursor = true

		Bus.Subscribe("endpoints:get", listEndpoints)

		v.Title = " Endpoints "
	}

	return nil
}

func listByBus() {

}

func listEndpoints(g *gocui.Gui, v *gocui.View, uri string) error {
	list, err := model.GetEndpoints(uri)

	if err != nil {
		return err
	}

	fmt.Fprintln(v, list)

	// if len(list) == 0 {
	// 	fmt.Fprintln(v, c.Red("not found uri"))
	// }

	for _, item := range list {
		fmt.Fprintln(v, "", c.Bold(getMethod(item["method"])), c.Gray(1, item["path"]))
	}

	return nil
}

var views = []string{"create-endpoint", "method", "path"}
var inputs = []string{"method", "path"}
var active = 0

func createView(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("")

	maxX, maxY := g.Size()

	x, x1 := maxX/5, maxX-(maxX/5)
	y, y1 := maxY/3, (maxY/2)-3

	if v, err := g.SetView("create-endpoint", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		uView, err := g.View("uri")

		if err != nil {
			return err
		}

		line, _ := uView.Line(0)

		if len(line) == 0 {
			v.Title = " Back to select uri "
			cmd.SetActive(g, "create-endpoint")

			fmt.Fprintln(v, "Provide uri after create a endpoint.")

			return nil
		}

		if err := elements.Input("method", x+2, y+1, x1-2, y+3)(g, v); err != nil {
			return err
		}

		if err := elements.Input("path", x+2, y+4, x1-2, y+6)(g, v); err != nil {
			return err
		}

		v.Title = " Create endpoint "
		if _, err := cmd.SetActive(g, "method"); err != nil {
			return err
		}
	}

	return nil
}

func createEndpoint(g *gocui.Gui, v *gocui.View) error {
	test := []string{"method", "path", "uri"}
	values := map[string]string{
		"method": "",
		"path":   "",
		"uri":    "",
	}

	for _, name := range test {

		key, err := g.View(name)
		if err != nil {
			return err
		}

		value, _ := key.Line(0)
		values[name] = value
	}

	if err := model.CreateEndpoint(values); err != nil {
		return err
	}

	if err := cmd.CloseList(views, "endpoints")(g, v); err != nil {
		return err
	}

	return nil
}

func eBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("endpoints", 'j', gocui.ModNone, cmd.CursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("endpoints", 'k', gocui.ModNone, cmd.CursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("endpoints", 'c', gocui.ModNone, createView); err != nil {
		return err
	}

	if err := g.SetKeybinding("endpoints", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}

	if err := g.SetKeybinding("endpoints", 'c', gocui.ModNone, createView); err != nil {
		return err
	}

	if err := g.SetKeybinding("create-endpoint", gocui.KeyCtrlC, gocui.ModNone, cmd.Close("create-endpoint", "endpoints")); err != nil {
		return err
	}

	for _, name := range inputs {
		if err := g.SetKeybinding(name, gocui.KeyCtrlC, gocui.ModNone, cmd.CloseList(views, "endpoints")); err != nil {
			return err
		}
		if err := g.SetKeybinding(name, gocui.KeyTab, gocui.ModNone, cmd.Navigate("next", inputs, &active)); err != nil {
			return err
		}
	}
	if err := g.SetKeybinding("method", gocui.KeyEnter, gocui.ModNone, cmd.Navigate("next", inputs, &active)); err != nil {
		return err
	}

	if err := g.SetKeybinding("path", gocui.KeyEnter, gocui.ModNone, createEndpoint); err != nil {
		return err
	}

	return nil
}

func getMethod(method string) c.Value {
	var selected c.Value
	switch method {
	case "GET":
		selected = c.Green(method)
	case "POST":
		selected = c.Yellow(method)
	case "PATCH":
		selected = c.Brown(method)
	case "PUT":
		selected = c.Brown(method)
	case "DELETE":
		selected = c.Brown(method)
	default:
		selected = c.Green(method)
	}

	return selected
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	// var l string
	// var err error

	// _, cy := v.Cursor()
	// if l, err = v.Line(cy); err != nil {
	// 	l = ""
	// }

	return nil
}
