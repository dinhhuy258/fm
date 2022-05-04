package command

import (
	"strconv"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/gocui"
)

func ToggleSelection(app IApp, _ ...interface{}) error {
	node := fs.GetFileManager().GetNodeAtIdx(app.GetFocusIdx())
	path := node.AbsolutePath

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

	fileManager := fs.GetFileManager()
	fileManager.Reload()

	numberOfFiles := fileManager.GetVisibleNodesSize()
	app.SetNumberOfFiles(numberOfFiles)
	title := (" " + fileManager.GetCurrentPath() + " (" + strconv.Itoa(numberOfFiles) + ") ")
	gui.GetGui().SetMainTitle(title)
	gui.GetGui().UpdateSortAndFilter()

	gui.GetGui().RenderDir(
		fs.GetFileManager().GetVisibleNodes(),
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
	fileManager := fs.GetFileManager()
	currentNode := fileManager.GetNodeAtIdx(app.GetFocusIdx())

	focus := currentNode.AbsolutePath
	if len(params) == 1 {
		focus, _ = params[0].(string)
	}

	ChangeDirectory(app, fileManager.GetCurrentPath(), false, &focus)

	return nil
}

func refreshSelections(app IApp) {
	gui.GetGui().RenderSelections(app.GetSelections())

	gui.GetGui().RenderDir(
		fs.GetFileManager().GetVisibleNodes(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)
}

func Quit(_ IApp, _ ...interface{}) error {
	return gocui.ErrQuit
}
