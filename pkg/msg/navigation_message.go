package msg

import (
	"fmt"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

// FocusByIndex focus entry by index
func FocusByIndex(app IApp, params ...string) {
	index, err := strconv.Atoi(params[0])
	if err != nil {
		logController, _ := app.GetController(controller.Log).(*controller.LogController)
		logController.SetLog(view.Error, fmt.Sprintf("Invalid index: %v", err))
		logController.UpdateView()
	}

	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	explorerController.FocusByIndex(index)
	explorerController.UpdateView()
}

// FocusFirst focus first entry
func FocusFirst(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusFirst()
	explorerController.UpdateView()
}

// FocusLast focus last entry
func FocusLast(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusLast()
	explorerController.UpdateView()
}

// FocusNext focus next entry
func FocusNext(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusNext()
	explorerController.UpdateView()
}

// FocusPrevious focus previous entry
func FocusPrevious(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.FocusPrevious()
	explorerController.UpdateView()
}

// FocusPath focus entry with path
func FocusPath(app IApp, params ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	path := params[0]
	if !fs.IsPathExists(path) {
		logController.SetLog(view.Error, "Path does not exist: "+path)
		logController.UpdateView()

		return
	}

	explorerController.FocusPath(path)
	explorerController.UpdateView()
}

// Enter directory
func Enter(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()
	if entry == nil {
		return
	}

	if entry.IsDirectory() || entry.IsSymlink() {
		explorerController.LoadDirectory(entry.GetPath(), optional.NewEmpty[string]())
		explorerController.UpdateView()
	}
}

// Back to parent directory
func Back(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	dir := fs.Dir(explorerController.GetPath())
	if dir == "." {
		// If folder has no parent directory then do nothing
		return
	}

	loadDirectory(app, dir, optional.New(explorerController.GetPath()))
}

// ChangeDirectory change directory
func ChangeDirectory(app IApp, params ...string) {
	directory := params[0]

	loadDirectory(app, directory, optional.NewEmpty[string]())
}

// loadDirectory load directory
func loadDirectory(app IApp, path string, focusPath optional.Optional[string]) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	explorerController.LoadDirectory(path, focusPath)
	explorerController.UpdateView()
}
