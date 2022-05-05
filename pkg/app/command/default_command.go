package command

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

func ToggleSelection(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()

	entry := fileExplorer.GetEntry(app.GetFocus())
	path := entry.GetPath()

	app.ToggleSelection(path)

	app.RenderSelections()
	app.RenderEntries()

	return nil
}

func ToggleHidden(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	appGui.UpdateSortAndFilter()

	entry := fileExplorer.GetEntry(app.GetFocus())
	LoadDirectory(app, fileExplorer.GetPath(), false, entry.GetPath())

	return nil
}

func ClearSelection(app IApp, _ ...interface{}) error {
	app.ClearSelections()

	app.RenderSelections()
	app.RenderEntries()

	return nil
}

func SwitchMode(app IApp, params ...interface{}) error {
	return app.PushMode(params[0].(string))
}

func PopMode(app IApp, _ ...interface{}) error {
	return app.PopMode()
}

func Refresh(app IApp, params ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()

	entry := fileExplorer.GetEntry(app.GetFocus())
	LoadDirectory(app, fileExplorer.GetPath(), false, entry.GetPath())

	return nil
}

func Quit(_ IApp, _ ...interface{}) error {
	appGui := gui.GetGui()

	return appGui.Quit()
}
