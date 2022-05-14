package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

func FocusFirst(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()

	appGui.GetControllers().Explorer.FocusFirst()

	return nil
}

func FocusNext(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()

	appGui.GetControllers().Explorer.FocusNext()

	return nil
}

func FocusPrevious(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()

	appGui.GetControllers().Explorer.FocusPrevious()

	return nil
}

func FocusPath(app IApp, params ...interface{}) error {
	appGui := gui.GetGui()

	//TODO Verify path
	path, _ := params[0].(string)
	appGui.GetControllers().Explorer.FocusPath(path)

	return nil
}

func Enter(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()
	explorerController := appGui.GetControllers().Explorer

	entry := explorerController.GetEntry(explorerController.GetFocus())
	if entry == nil {
		return nil
	}

	if entry.IsDirectory() {
		explorerController.LoadDirectory(entry.GetPath(), "")
	}

	return nil
}

func Back(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()

	dir := fileExplorer.Dir()
	if dir == "." {
		// If folder has no parent directory then do nothing
		return nil
	}

	LoadDirectory(app, dir, fileExplorer.GetPath())

	return nil
}

func ChangeDirectory(app IApp, params ...interface{}) error {
	directory, _ := params[0].(string)

	LoadDirectory(app, directory, "")

	return nil
}

//TODO: Remove?
func LoadDirectory(app IApp, path string, focusPath string) {
	appGui := gui.GetGui()

	appGui.GetControllers().Explorer.LoadDirectory(path, "")
}
