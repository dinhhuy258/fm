package command

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

func ToggleSelection(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()
	explorerController := appGui.GetControllers().Explorer
	selectionController := appGui.GetControllers().Sellection

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
	appGui := gui.GetGui()
	explorerController := appGui.GetControllers().Explorer


	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	entry := explorerController.GetCurrentEntry()
	LoadDirectory(app, fs.GetFileExplorer().GetPath(), entry.GetPath())

	return nil
}

func ClearSelection(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()
	explorerController := appGui.GetControllers().Explorer
	selectionController := appGui.GetControllers().Sellection

	selectionController.ClearSelections()
	explorerController.UpdateView()

	return nil
}

func SwitchMode(app IApp, params ...interface{}) error {
	return app.PushMode(params[0].(string))
}

func PopMode(app IApp, _ ...interface{}) error {
	return app.PopMode()
}

func Refresh(app IApp, params ...interface{}) error {
	appGui := gui.GetGui()
	explorerController := appGui.GetControllers().Explorer
	fileExplorer := fs.GetFileExplorer()

	entry := explorerController.GetCurrentEntry()
	LoadDirectory(app, fileExplorer.GetPath(), entry.GetPath())

	return nil
}

func Quit(_ IApp, _ ...interface{}) error {
	appGui := gui.GetGui()

	return appGui.Quit()
}
