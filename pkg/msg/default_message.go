package msg

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/optional"
)

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

	app.SwitchMode(mode)
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
