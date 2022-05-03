package command

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/gocui"
)

func ToggleSelection(app IApp, _ ...interface{}) error {
	path := fs.GetFileManager().Dir.VisibleNodes[app.State().FocusIdx].AbsolutePath

	if _, hasPath := app.State().Selections[path]; hasPath {
		delete(app.State().Selections, path)
	} else {
		app.State().Selections[path] = struct{}{}
	}

	refreshSelections(app)

	return nil
}

func ToggleHidden(app IApp, _ ...interface{}) error {
	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	fs.GetFileManager().Dir.Reload()

	numberOfFiles := len(fs.GetFileManager().Dir.VisibleNodes)
	app.State().NumberOfFiles = numberOfFiles
	gui.GetGui().Views.Main.SetTitle(fs.GetFileManager().Dir.Path, numberOfFiles)
	gui.GetGui().Views.SortAndFilter.SetSortAndFilter()

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		app.State().Selections,
		app.State().FocusIdx,
	)

	return nil
}

func ClearSelection(app IApp, _ ...interface{}) error {
	app.State().Selections = make(map[string]struct{})

	refreshSelections(app)

	return nil
}

func SwitchMode(app IApp, params ...interface{}) error {
	return app.PushMode(params[0].(string))
}

func PopMode(app IApp, _ ...interface{}) error {
	return app.PopMode()
}

func Refresh(app IApp, params ...interface{}) error {
	currentNode := fs.GetFileManager().Dir.VisibleNodes[app.State().FocusIdx]

	focus := currentNode.AbsolutePath
	if len(params) == 1 {
		focus = params[0].(string)
	}

	ChangeDirectory(app, fs.GetFileManager().Dir.Path, false, &focus)

	return nil
}

func refreshSelections(app IApp) {
	gui.GetGui().Views.Selection.SetTitle(len(app.State().Selections))
	gui.GetGui().Views.Selection.RenderSelections(app.State().Selections)

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		app.State().Selections,
		app.State().FocusIdx,
	)
}

func Quit(_ IApp, _ ...interface{}) error {
	return gocui.ErrQuit
}
