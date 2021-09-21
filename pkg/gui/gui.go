package gui

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/gocui"
)

type Gui struct {
	g              *gocui.Gui
	Views          *view.Views
	onViewsCreated func()
}

func NewGui(onViewsCreated func()) (*Gui, error) {
	gui := &Gui{
		onViewsCreated: onViewsCreated,
	}

	return gui, nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}

	gui.g = g
	gui.g.Cursor = true
	gui.g.InputEsc = true

	defer gui.g.Close()

	gui.g.SetManager(gocui.ManagerFunc(gui.layout))

	err = gui.g.MainLoop()

	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}

	return nil
}

func (gui *Gui) SetOnKeyFunc(onKey func(string) error) {
	gui.g.SetOnKeyFunc(onKey)
}
