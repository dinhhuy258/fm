package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/key"
)

// ToggleSelection is a message that toggles the selection of the current entry
func ToggleSelection(app IApp, _ key.Key, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Selection).(*controller.SelectionController)

	entry := explorerController.GetCurrentEntry()
	if entry == nil {
		return
	}

	path := entry.GetPath()

	selectionController.ToggleSelection(path)

	selectionController.UpdateView()
	explorerController.UpdateView()
}

// ClearSelection is a message that clears the selection
func ClearSelection(app IApp, _ key.Key, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Selection).(*controller.SelectionController)

	selectionController.ClearSelections()

	selectionController.UpdateView()
	explorerController.UpdateView()
}

// SelectAll is a message that select all visible entries
func SelectAll(app IApp, _ key.Key, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Selection).(*controller.SelectionController)

	entries := explorerController.GetEntries()
	for _, entry := range entries {
		path := entry.GetPath()
		selectionController.SelectPath(path)
	}

	selectionController.UpdateView()
	explorerController.UpdateView()
}

// ToggleSelectionByPath is a message that toggle selection by file path.
func ToggleSelectionByPath(app IApp, _ key.Key, params ...string) {
	path := params[0]

	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Selection).(*controller.SelectionController)

	selectionController.ToggleSelection(path)

	selectionController.UpdateView()
	explorerController.UpdateView()
}
