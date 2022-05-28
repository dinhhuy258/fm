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

func NewGui() (*Gui, error) {
	gui := &Gui{}

	g, err := gocui.NewGui(gocui.OutputNormal, false, gocui.NORMAL, false, map[rune]string{})
	if err != nil {
		return nil, err
	}

	gui.g = g
	gui.g.Cursor = false
	gui.g.InputEsc = true
	gui.g.SetManager(gocui.ManagerFunc(gui.layout))

	gui.views = view.CreateViews(gui.g)

	if _, err := gui.g.SetCurrentView(gui.views.Input.Name()); err != nil {
		return nil, err
	}

	gui.controllers = controller.CreateControllers(gui.g, gui.views)

	if err := gui.layout(gui.g); err != nil {
		return nil, err
	}

	return gui, nil
}

func (gui *Gui) Run() error {
	defer gui.g.Close()

	err := gui.g.MainLoop()
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

// Render re-render all views
func (gui *Gui) Render() {
	gui.OnUIThread(func() error { return nil })
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
