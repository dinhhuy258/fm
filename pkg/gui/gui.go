package gui

import (
	"errors"
	"sync"

	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/gocui"
)

type Gui struct {
	g              *gocui.Gui
	views          *view.Views
	controllers    *controller.Controllers
	onViewsCreated func()
}

func (gui *Gui) GetControllers() *controller.Controllers {
	return gui.controllers
}

func (gui *Gui) SetConfirmation(ask string, onConfirm func(bool)) {
	gui.views.Confirm.SetConfirmation(ask, func(ans bool) {
		gui.views.Explorer.SetAsCurrentView()

		onConfirm(ans)
	})
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

	gui.views = view.CreateAllViews(gui.g)
	gui.controllers = controller.CreateAllControllers()
	gui.controllers.Explorer.SetView(gui.views.Explorer)
	gui.controllers.Help.SetView(gui.views.Help)
	gui.controllers.Sellection.SetView(gui.views.Selection)
	gui.controllers.Progress.SetView(gui.views.Progress)
	gui.controllers.Log.SetView(gui.views.Log)
	gui.controllers.Input.SetView(gui.views.Input)

	gui.layout(gui.g)
	gui.onViewsCreated()

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
