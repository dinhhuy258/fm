package msg

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func ToggleSelection(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)

	entry := explorerController.GetCurrentEntry()
	if entry == nil {
		return
	}

	path := entry.GetPath()

	selectionController.ToggleSelection(path)

	selectionController.UpdateView()
	explorerController.UpdateView()
}

func ToggleHidden(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	entry := explorerController.GetCurrentEntry()
	loadDirectory(app, explorerController.GetPath(), optional.New(entry.GetPath()))
}

func ClearSelection(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)

	selectionController.ClearSelections()

	selectionController.UpdateView()
	explorerController.UpdateView()
}

func SwitchMode(app IApp, params ...string) {
	mode := params[0]

	app.PushMode(mode)
}

func PopMode(app IApp, _ ...string) {
	app.PopMode()
}

func Refresh(app IApp, params ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()

	focusPath := optional.NewEmpty[string]()
	if entry != nil {
		focusPath = optional.New(entry.GetPath())
	}

	loadDirectory(app, explorerController.GetPath(), focusPath)
}

func Quit(app IApp, _ ...string) {
	app.Quit()
}
