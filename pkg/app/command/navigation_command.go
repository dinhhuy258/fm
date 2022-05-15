package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
)

func FocusFirst(app IApp, _ ...interface{}) error {
	appGui := app.GetGui()

	appGui.GetControllers().Explorer.FocusFirst()

	return nil
}

func FocusNext(app IApp, _ ...interface{}) error {
	appGui := app.GetGui()

	appGui.GetControllers().Explorer.FocusNext()

	return nil
}

func FocusPrevious(app IApp, _ ...interface{}) error {
	appGui := app.GetGui()

	appGui.GetControllers().Explorer.FocusPrevious()

	return nil
}

func FocusPath(app IApp, params ...interface{}) error {
	appGui := app.GetGui()

	// TODO Verify path
	path, _ := params[0].(string)
	appGui.GetControllers().Explorer.FocusPath(path)

	return nil
}

func Enter(app IApp, _ ...interface{}) error {
	appGui := app.GetGui()
	explorerController := appGui.GetControllers().Explorer

	entry := explorerController.GetCurrentEntry()
	if entry == nil {
		return nil
	}

	if entry.IsDirectory() {
		explorerController.LoadDirectory(entry.GetPath(), "")
	}

	return nil
}

func Back(app IApp, _ ...interface{}) error {
	explorerController := app.GetGui().GetControllers().Explorer

	dir := fs.Dir(explorerController.GetPath())
	if dir == "." {
		// If folder has no parent directory then do nothing
		return nil
	}

	loadDirectory(app, dir, explorerController.GetPath())

	return nil
}

func ChangeDirectory(app IApp, params ...interface{}) error {
	directory, _ := params[0].(string)

	loadDirectory(app, directory, "")

	return nil
}

func loadDirectory(app IApp, path string, focusPath string) {
	appGui := app.GetGui()

	appGui.GetControllers().Explorer.LoadDirectory(path, "")
}
