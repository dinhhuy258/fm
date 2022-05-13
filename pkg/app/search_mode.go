package app

import (
	"strings"

	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/fs"
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

		if searchInput != "" {
			fileExplorer := fs.GetFileExplorer()

			entries := fileExplorer.GetEntries()
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
