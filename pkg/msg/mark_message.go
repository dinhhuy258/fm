package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/key"
)

func MarkSave(app IApp, params ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	marController, _ := app.GetController(controller.Mark).(*controller.MarkController)

	entry := explorerController.GetCurrentEntry()
	if entry == nil {
		return
	}

	k := key.GetKeyDisplay(app.GetPressedKey())

	marController.SaveMark(k, entry.GetPath())
}

func MarkLoad(app IApp, params ...string) {
	marController, _ := app.GetController(controller.Mark).(*controller.MarkController)

	k := key.GetKeyDisplay(app.GetPressedKey())

	path := marController.LoadMark(k)

	path.IfPresent(func(p *string) {
		FocusPath(app, *p)
	})
}
