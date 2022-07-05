package views

import (
	"github.com/asaskevich/EventBus"
	"github.com/jroimartin/gocui"
)

var Bus = EventBus.New()

func Layouts(g *gocui.Gui) error {
	if err := Uri(g); err != nil {
		return err
	}

	if err := Endpoints(g); err != nil {
		return err
	}

	if err := Content(g); err != nil {
		return err
	}

	return nil
}
