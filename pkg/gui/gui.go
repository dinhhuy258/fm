package gui

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/row"
	"github.com/dinhhuy258/gocui"
)

type Views struct {
	MainHeader    *gocui.View
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
	MainRow       *row.MainRow
	GuiLoadedChan chan struct{}
}

func NewGui() (*Gui, error) {
	gui := &Gui{
		ViewsSetup:    false,
		GuiLoadedChan: make(chan struct{}, 1),
		MainRow:       row.NewMainRow(),
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
