package views

import (
	"fmt"

	"github.com/gandarfh/httui/cmd"
	"github.com/gandarfh/httui/db"
	"github.com/jroimartin/gocui"
	c "github.com/logrusorgru/aurora/v3"
)

func navigate(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
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

func cursorUp(g *gocui.Gui, v *gocui.View) error {
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

func getLine(g *gocui.Gui, v *gocui.View) error {
	// var l string
	// var err error

	// _, cy := v.Cursor()
	// if l, err = v.Line(cy); err != nil {
	// 	l = ""
	// }

	return nil
}

func saveNewEndpoint(g *gocui.Gui, v *gocui.View) error {
	db, err := db.Conn()

	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`INSERT INTO uris (uri) VALUES (?)`)

	if err != nil {
		return err
	}

	value, err := v.Line(0)

	if err != nil {
		return err
	}

	stmt.Exec(value)

	db.Close()
	closeCreateUri(g, v)
	return nil
}

func newEndpoint(g *gocui.Gui, v *gocui.View) error {
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

	if err := g.SetKeybinding("create-uri", gocui.KeyEnter, gocui.ModNone, saveNewEndpoint); err != nil {
		return err
	}

	if err := g.SetKeybinding("create-uri", gocui.KeyCtrlC, gocui.ModNone, closeCreateUri); err != nil {
		return err
	}

	return nil
}

func closeCreateUri(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("create-uri"); err != nil {
		return err
	}

	if err := cmd.Bindings(g); err != nil {
		return err
	}

	if _, err := cmd.SetActive(g, "endpoints"); err != nil {
		return err
	}

	return nil
}

func Endpoints(g *gocui.Gui) error {
	_, maxY := g.Size()

	if v, err := g.SetView("endpoints", 4, 6, 30, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if _, err = cmd.SetActive(g, "endpoints"); err != nil {
			return err
		}

		if err := g.SetKeybinding("endpoints", 'j', gocui.ModNone, cursorDown); err != nil {
			return err
		}

		if err := g.SetKeybinding("endpoints", 'k', gocui.ModNone, cursorUp); err != nil {
			return err
		}

		if err := g.SetKeybinding("endpoints", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
			return err
		}

		if err := g.SetKeybinding("endpoints", 'c', gocui.ModNone, newEndpoint); err != nil {
			return err
		}

		db, err := db.Conn()

		if err != nil {
			return err
		}

		db.Close()

		v.Autoscroll = false

		g.Cursor = true

		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		fmt.Fprintln(v, "users")
		fmt.Fprintln(v, "", c.Bold(c.Green("GET")), c.Gray(1, "api/users"))
		fmt.Fprintln(v, "", c.Bold(c.Green("GET")), c.Gray(1, "api/users/[:id]"))
		fmt.Fprintln(v, "", c.Bold(c.Yellow("POST")), c.Gray(1, "api/users"))
		fmt.Fprintln(v, "", c.Bold(c.Brown("PETCH")), c.Gray(1, "api/users/[:id]"))
		fmt.Fprintln(v, "", c.Bold(c.Yellow("PUT")), c.Gray(1, "api/users/[:id]"))
		fmt.Fprintln(v, "", c.Bold(c.Red("DELETE")), c.Gray(1, "api/users/[:id]"))

		v.Title = " Endpoints "

	}

	return nil
}
