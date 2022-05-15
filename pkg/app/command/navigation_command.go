package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func FocusFirst(app IApp, _ ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusFirst()

	return nil
}

func FocusNext(app IApp, _ ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusNext()

	return nil
}

func FocusPrevious(app IApp, _ ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusPrevious()

	return nil
}

func FocusPath(app IApp, params ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	// TODO Verify path
	path, _ := params[0].(string)
	explorerController.FocusPath(path)

	return nil
}

func Enter(app IApp, _ ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()
	if entry == nil {
		return nil
	}

	if entry.IsDirectory() {
		explorerController.LoadDirectory(entry.GetPath(), optional.NewEmptyOptional[string]())
	}

	return nil
}

func Back(app IApp, _ ...interface{}) error {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	dir := fs.Dir(explorerController.GetPath())
	if dir == "." {
		// If folder has no parent directory then do nothing
		return nil
	}

	loadDirectory(app, dir, optional.NewOptional(explorerController.GetPath()))

	return nil
}

func ChangeDirectory(app IApp, params ...interface{}) error {
	directory, _ := params[0].(string)

	loadDirectory(app, directory, optional.NewEmptyOptional[string]())

	return nil
}

func loadDirectory(app IApp, path string, focusPath optional.Optional[string]) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.LoadDirectory(path, focusPath)
}
