package msg

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/optional"
)

// ToggleSelection is a message that toggles the selection of the current entry
func ToggleSelection(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)

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
func ClearSelection(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)

	selectionController.ClearSelections()

	selectionController.UpdateView()
	explorerController.UpdateView()
}

// ToggleHidden is a message that toggles the hidden configuration
func ToggleHidden(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	config.AppConfig.General.ShowHidden = !config.AppConfig.General.ShowHidden

	entry := explorerController.GetCurrentEntry()
	loadDirectory(app, explorerController.GetPath(), optional.New(entry.GetPath()))
}

// SwitchMode is a message that switches the mode of the application
func SwitchMode(app IApp, params ...string) {
	mode := params[0]

	app.PushMode(mode)
}

// PopMode is a message that pops the current mode
func PopMode(app IApp, _ ...string) {
	app.PopMode()
}

// Refresh is a message that refreshes the current directory
func Refresh(app IApp, params ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()

	focusPath := optional.NewEmpty[string]()
	if entry != nil {
		focusPath = optional.New(entry.GetPath())
	}

	loadDirectory(app, explorerController.GetPath(), focusPath)
}

// Quit is a message that quits the application
func Quit(app IApp, _ ...string) {
	app.Quit()
}
