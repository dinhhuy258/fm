package app

import (
	"strings"

	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
)

type SearchMode struct {
	*Mode
}

func (*SearchMode) GetName() string {
	return "search"
}

func (m *SearchMode) OnModeStarted(app *App) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	inputController.SetInput(controller.InputText, "search", func(searchInput string) {
		command.PopMode(app)

		if searchInput != "" {
			entries := explorerController.GetEntries()
			for _, entry := range entries {
				if strings.Contains(strings.ToLower(entry.GetName()), strings.ToLower(searchInput)) {
					command.FocusPath(app, entry.GetPath())

					return
				}
			}
		}
	})
}

func createSearchMode() *SearchMode {
	return &SearchMode{
		&Mode{
			keyBindings: &KeyBindings{
				OnKeys: map[string]*Action{
					"esc": {
						Help: "cancel",
						Commands: []*command.Command{
							{
								Func: command.PopMode,
							},
						},
					},
				},
			},
		},
	}
}
