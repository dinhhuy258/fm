package gui

import (
	"errors"
	"sync"

	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/gocui"
)

type Gui struct {
	g              *gocui.Gui
	Views          *view.Views
	onViewsCreated func()
}

var (
	gui                   *Gui
	guiInitializationOnce sync.Once
)

func InitGui(onViewsCreated func()) {
	// Make sure only one gui is created
	guiInitializationOnce.Do(func() {
		gui = &Gui{
			onViewsCreated: onViewsCreated,
		}
	})
}

func GetGui() *Gui {
	if gui == nil {
		panic("gui is not initialized")
	}

	return gui
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
