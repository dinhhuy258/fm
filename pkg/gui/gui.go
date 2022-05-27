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

func NewGui() *Gui {
	return &Gui{}
}

func (gui *Gui) Run(onGuiReady func()) error {
	g, err := gocui.NewGui(gocui.OutputNormal, false, gocui.NORMAL, false, map[rune]string{})
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

func (gui *Gui) SetOnKeyFunc(onKey func(key gocui.Key, ch rune, mod gocui.Modifier) error) {
	gui.g.SetOnKeyFunc(onKey)
}

func (gui *Gui) Quit() {
	gui.g.Update(func(g *gocui.Gui) error {
		return gocui.ErrQuit
	})
}

func (gui *Gui) OnUIThread(f func() error) {
	gui.g.Update(func(*gocui.Gui) error {
		return f()
	})
}

func (gui *Gui) Suspend() error {
	return gui.g.Suspend()
}

func (gui *Gui) Resume() error {
	return gui.g.Resume()
}

func (gui *Gui) GetController(controllerType controller.Type) controller.IController {
	return gui.controllers.GetController(controllerType)
}
