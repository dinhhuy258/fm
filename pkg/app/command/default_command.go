package command

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func ToggleSelection(app IApp, _ ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)

	entry := explorerController.GetCurrentEntry()
	if entry == nil {
		return nil
	}

	path := entry.GetPath()

	selectionController.ToggleSelection(path)

	explorerController.UpdateView()

	return nil
}

func ToggleHidden(app IApp, _ ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	entry := explorerController.GetCurrentEntry()
	loadDirectory(app, explorerController.GetPath(), optional.NewOptional(entry.GetPath()))

	return nil
}

func ClearSelection(app IApp, _ ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)

	selectionController.ClearSelections()
	explorerController.UpdateView()

	return nil
}

func SwitchMode(app IApp, params ...interface{}) error {
	mode, _ := params[0].(string)

	return app.PushMode(mode)
}

func PopMode(app IApp, _ ...interface{}) error {
	return app.PopMode()
}

func Refresh(app IApp, params ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()
	loadDirectory(app, explorerController.GetPath(), optional.NewOptional(entry.GetPath()))

	return nil
}

func Quit(app IApp, _ ...interface{}) error {
	return app.Quit()
}
