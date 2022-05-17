package command

import "github.com/dinhhuy258/fm/pkg/gui/controller"

func MarkSave(app IApp, params ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)

	entry := explorerController.GetCurrentEntry()
	key, _ := params[0].(string)
	app.MarkSave(key, entry.GetPath())
}

func MarkLoad(app IApp, params ...interface{}) {
	key, _ := params[0].(string)

	if path, hasKey := app.MarkLoad(key); hasKey {
		FocusPath(app, path)
	}
}
