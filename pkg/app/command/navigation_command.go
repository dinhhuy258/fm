package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func FocusFirst(app IApp, _ ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusFirst()
	explorerController.UpdateView()
}

func FocusNext(app IApp, _ ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusNext()
	explorerController.UpdateView()
}

func FocusPrevious(app IApp, _ ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusPrevious()
	explorerController.UpdateView()
}

func FocusPath(app IApp, params ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	// TODO Verify path
	path, _ := params[0].(string)
	explorerController.FocusPath(path)
	explorerController.UpdateView()
}

func Enter(app IApp, _ ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()
	if entry == nil {
		return
	}

	if entry.IsDirectory() {
		explorerController.LoadDirectory(entry.GetPath(), optional.NewEmptyOptional[string]())
		explorerController.UpdateView()
	}
}

func Back(app IApp, _ ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	dir := fs.Dir(explorerController.GetPath())
	if dir == "." {
		// If folder has no parent directory then do nothing
		return
	}

	loadDirectory(app, dir, optional.NewOptional(explorerController.GetPath()))
}

func ChangeDirectory(app IApp, params ...interface{}) {
	directory, _ := params[0].(string)

	loadDirectory(app, directory, optional.NewEmptyOptional[string]())
}

// TODO: Considering remove this method and use ChangeDirectory instead
func loadDirectory(app IApp, path string, focusPath optional.Optional[string]) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.LoadDirectory(path, focusPath)
	explorerController.UpdateView()
}
