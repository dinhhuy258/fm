package command

import (
	"strconv"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/gocui"
)

func ToggleSelection(app IApp, _ ...interface{}) error {
	path := fs.GetFileManager().Dir.VisibleNodes[app.GetFocusIdx()].AbsolutePath

	if app.HasSelection(path) {
		app.DeleteSelection(path)
	} else {
		app.AddSelection(path)
	}

	refreshSelections(app)

	return nil
}

func ToggleHidden(app IApp, _ ...interface{}) error {
	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	fs.GetFileManager().Dir.Reload()

	numberOfFiles := len(fs.GetFileManager().Dir.VisibleNodes)
	app.SetNumberOfFiles(numberOfFiles)
	title := (" " + fs.GetFileManager().Dir.Path + " (" + strconv.Itoa(numberOfFiles) + ") ")
	gui.GetGui().SetMainTitle(title)
	gui.GetGui().UpdateSortAndFilter()

	gui.GetGui().RenderDir(
		fs.GetFileManager().Dir,
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func ClearSelection(app IApp, _ ...interface{}) error {
	app.ClearSelections()

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
	currentNode := fs.GetFileManager().Dir.VisibleNodes[app.GetFocusIdx()]

	focus := currentNode.AbsolutePath
	if len(params) == 1 {
		focus, _ = params[0].(string)
	}

	ChangeDirectory(app, fs.GetFileManager().Dir.Path, false, &focus)

	return nil
}

func refreshSelections(app IApp) {
	gui.GetGui().RenderSelections(app.GetSelections())

	gui.GetGui().RenderDir(
		fs.GetFileManager().Dir,
		app.GetSelections(),
		app.GetFocusIdx(),
	)
}

func Quit(_ IApp, _ ...interface{}) error {
	return gocui.ErrQuit
}
