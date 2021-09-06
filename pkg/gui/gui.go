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

type State struct {
	Main *MainPanelState
}

type Gui struct {
	g             *gocui.Gui
	ViewsSetup    bool
	Views         Views
	State         *State
	GuiLoadedChan chan struct{}
}

func NewGui() (*Gui, error) {
	gui := &Gui{
		ViewsSetup: false,
		State: &State{
			Main: &MainPanelState{
				SelectedIdx:   0,
				NumberOfFiles: 0,
			},
		},
		GuiLoadedChan: make(chan struct{}, 1),
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

	if err = gui.setKeyBindings(); err != nil {
		return err
	}

	err = gui.g.MainLoop()

	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}

	return nil
}
