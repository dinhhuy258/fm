package gui

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/gocui"
)

type Gui struct {
	g           *gocui.Gui
	views       *view.Views
	controllers *controller.Controllers
}

func (gui *Gui) GetControllers() *controller.Controllers {
	return gui.controllers
}

func NewGui() *Gui {
	return &Gui{}
}

func (gui *Gui) Run(onGuiReady func()) error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}

	gui.g = g
	gui.g.Cursor = true
	gui.g.InputEsc = true

	defer gui.g.Close()

	gui.g.SetManager(gocui.ManagerFunc(gui.layout))

	gui.views = view.CreateAllViews(gui.g)
	gui.controllers = controller.CreateAllControllers(gui.views)

	if err := gui.layout(gui.g); err != nil {
		return err
	}

	onGuiReady()

	err = gui.g.MainLoop()

	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}

	return nil
}

func (gui *Gui) SetOnKeyFunc(onKey func(string) error) {
	gui.g.SetOnKeyFunc(onKey)
}

func (gui *Gui) Quit() error {
	return gocui.ErrQuit
}
