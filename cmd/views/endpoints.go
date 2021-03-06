package views

import (
	"fmt"
	"strings"

	"github.com/gandarfh/httui/cmd"
	"github.com/gandarfh/httui/cmd/elements"
	"github.com/gandarfh/httui/cmd/model"
	"github.com/jroimartin/gocui"
	c "github.com/logrusorgru/aurora/v3"
)

func Endpoints(g *gocui.Gui, config *cmd.Config) error {
	_, maxY := g.Size()

	x, x1 := 4, 36
	y, y1 := 5, maxY-4

	if v, err := g.SetView("endpoints", x, y, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := eBindings(g, config); err != nil {
			return err
		}

		v.Autoscroll = true
		g.Cursor = true

		listEndpoints(config)(g, v, &config.Default)
		cmd.Bus.Subscribe("endpoints:get", listEndpoints(config))

		v.Title = " Endpoints "
	}

	return nil
}

func listEndpoints(config *cmd.Config) func(g *gocui.Gui, v *gocui.View, uri *string) error {
	return func(g *gocui.Gui, v *gocui.View, uri *string) error {
		list := config.GetUri(uri)

		config.SetDefaultUri(uri)

		if len(*list.ListEndpoints()) == 0 {
			fmt.Fprintln(v, c.Red("not found uri"))
		}

		for _, item := range *list.Endpoints {
			fmt.Fprintln(v, "", c.Bold(getMethod(item.Method)), c.Gray(1, item.Path))
		}

		return nil

	}
}

var views = []string{"create-endpoint", "method", "path"}
var inputs = []string{"method", "path"}
var active = 0

func createView(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("")

	maxX, maxY := g.Size()

	x, x1 := maxX/3, maxX-(maxX/3)
	y, y1 := maxY/3, (maxY/3)+7

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

func createEndpoint(config *cmd.Config) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		tests := []string{"method", "path", "uri"}
		values := map[string]string{
			"method": "",
			"path":   "",
			"uri":    "",
		}

		for _, name := range tests {
			key, err := g.View(name)
			if err != nil {
				return err
			}

			value, _ := key.Line(0)
			values[name] = value
		}

		address := values["uri"]

		uri := config.GetUri(&address)

		endpoint := model.Endpoint{Path: values["path"], Method: strings.ToUpper(values["method"]), Headers: values["uri"]}

		err := config.CreateEndpoint(uri, &endpoint)

		if err != nil {
			return err
		}

		eView, err := g.View("endpoints")

		if err != nil {
			return err
		}

		eView.Clear()
		cmd.Bus.Publish("endpoints:get", g, eView, &uri.Alias)

		if err := cmd.CloseList(views, "endpoints")(g, v); err != nil {
			return err
		}

		return nil

	}
}

func eBindings(g *gocui.Gui, config *cmd.Config) error {
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

	if err := g.SetKeybinding("path", gocui.KeyEnter, gocui.ModNone, createEndpoint(config)); err != nil {
		return err
	}

	return nil
}

func getMethod(method string) c.Value {
	var selected c.Value

	switch strings.ToUpper(method) {
	case "GET":
		selected = c.Green(method)
	case "POST":
		selected = c.Yellow(method)
	case "PATCH":
		selected = c.Brown(method)
	case "PUT":
		selected = c.Brown(method)
	case "DELETE":
		selected = c.Red(method)
	default:
		selected = c.Green(method)
	}

	return selected
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	line, _ := cmd.GetLine(v)
	// vContent, _ := g.View("content")

	value := strings.Split(strings.TrimLeft(line, " "), " ")

	cmd.Bus.Publish("test", value[1])

	// vContent.Title = fmt.Sprintf(" %v ", value[1])

	return nil
}
