package app

import (
	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/gui"
)

type SearchMode struct {
	*Mode
}

func (*SearchMode) GetName() string {
	return "search"
}

func (m *SearchMode) OnModeStarted(app *App) {
	appGui := gui.GetGui()

	appGui.SetInput("search", func(searchInput string) {
		_ = command.PopMode(app)
		appGui.SetLogViewOnTop()
	})
}

func createSearchMode() *SearchMode {
	return &SearchMode{
		&Mode{
			KeyBindings: &KeyBindings{
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
