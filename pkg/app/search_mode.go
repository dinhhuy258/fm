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
	appGui := app.GetGui()
	explorerControler := appGui.GetControllers().Explorer
	inputController := appGui.GetControllers().Input

	inputController.SetInput(controller.Input, "search", func(searchInput string) {
		_ = command.PopMode(app)

		if searchInput != "" {
			entries := explorerControler.GetEntries()
			for _, entry := range entries {
				if strings.Contains(strings.ToLower(entry.GetName()), strings.ToLower(searchInput)) {
					_ = command.FocusPath(app, entry.GetPath())

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
