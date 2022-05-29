package msg

import (
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func FocusFirst(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusFirst()
	explorerController.UpdateView()
}

func FocusLast(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusLast()
	explorerController.UpdateView()
}

func FocusNext(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusNext()
	explorerController.UpdateView()
}

func FocusPrevious(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusPrevious()
	explorerController.UpdateView()
}

func FocusPath(app IApp, params ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	path := params[0]
	if !fs.IsDir(path) {
		logController.SetLog(view.Error, "Path is not a directory: "+path)

		return
	}

	explorerController.FocusPath(path)
	explorerController.UpdateView()
}

func Enter(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()
	if entry == nil {
		return
	}

	if entry.IsDirectory() {
		explorerController.LoadDirectory(entry.GetPath(), optional.NewEmpty[string]())
		explorerController.UpdateView()
	}
}

func Back(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	dir := fs.Dir(explorerController.GetPath())
	if dir == "." {
		// If folder has no parent directory then do nothing
		return
	}

	loadDirectory(app, dir, optional.New(explorerController.GetPath()))
}

func ChangeDirectory(app IApp, params ...string) {
	directory := params[0]

	loadDirectory(app, directory, optional.NewEmpty[string]())
}

// TODO: Considering remove this method and use ChangeDirectory instead
func loadDirectory(app IApp, path string, focusPath optional.Optional[string]) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.LoadDirectory(path, focusPath)
	explorerController.UpdateView()
}
