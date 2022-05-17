package command

import (
	"strings"

	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func Search(app IApp, _ ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	inputController.SetInput(controller.InputText, "search", optional.NewEmpty[string](),
		func(searchInput string) {
			if searchInput != "" {
				entries := explorerController.GetEntries()
				for _, entry := range entries {
					if strings.Contains(strings.ToLower(entry.GetName()), strings.ToLower(searchInput)) {
						FocusPath(app, entry.GetPath())

						return
					}
				}
			}
		})
	inputController.UpdateView()
}
