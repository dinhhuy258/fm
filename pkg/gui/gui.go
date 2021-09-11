package gui

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/gocui"
)

type Gui struct {
	g                *gocui.Gui
	Views            *view.Views
	ViewsCreatedChan chan struct{}
}

func NewGui() (*Gui, error) {
	gui := &Gui{
		ViewsCreatedChan: make(chan struct{}, 1),
	}

	return gui, nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}

	gui.g = g

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
