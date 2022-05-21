package command

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/key"
)

func MarkSave(app IApp, params ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()
	k := key.GetKeyDisplay(app.GetPressedKey())
	app.MarkSave(k, entry.GetPath())
}

func MarkLoad(app IApp, params ...interface{}) {
	k := key.GetKeyDisplay(app.GetPressedKey())

	if path, hasKey := app.MarkLoad(k); hasKey {
		FocusPath(app, path)
	}
}
