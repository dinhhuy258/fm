package message

import (
	"strings"

	"github.com/dinhhuy258/fm/pkg/gui/controller"
)

func SearchFromInput(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)
	filterValue := inputController.GetInputBuffer()

	if filterValue != "" {
		entries := explorerController.GetEntries()
		for _, entry := range entries {
			if strings.Contains(strings.ToLower(entry.GetName()), strings.ToLower(filterValue)) {
				FocusPath(app, entry.GetPath())

				return
			}
		}
	}
}
