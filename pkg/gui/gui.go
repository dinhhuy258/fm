package gui

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
)

type Gui struct {
	g *gocui.Gui
}

func NewGui() (*Gui, error) {
	gui := &Gui{}

	return gui, nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}

	gui.g = g

	defer gui.g.Close()

	gui.g.SetManagerFunc(layout)

	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	err = gui.g.MainLoop()

	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		fmt.Fprintln(v, "Hello world!")
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
