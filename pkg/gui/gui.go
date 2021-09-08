package gui

import (
	"errors"

	"github.com/jroimartin/gocui"
)

type Views struct {
	Main          *gocui.View
	Selection     *gocui.View
	HelpMenu      *gocui.View
	SortAndFilter *gocui.View
	InputAndLogs  *gocui.View
}

type Gui struct {
	g             *gocui.Gui
	ViewsSetup    bool
	Views         Views
	GuiLoadedChan chan struct{}
}

func NewGui() (*Gui, error) {
	gui := &Gui{
		ViewsSetup:    false,
		GuiLoadedChan: make(chan struct{}, 1),
	}

	return gui, nil
}

func (gui *Gui) Run(onKey func(string) error) error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}

	gui.g = g

	defer gui.g.Close()

	gui.g.SetManager(gocui.ManagerFunc(gui.layout))

	if err = gui.setKeyBindings(onKey); err != nil {
		return err
	}

	err = gui.g.MainLoop()

	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}

	return nil
}
