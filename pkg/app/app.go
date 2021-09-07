package app

import (
	"log"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

type App struct {
	Gui         *gui.Gui
	Mode        *Mode
	FileManager *fs.FileManager
}

// NewApp bootstrap a new application
func NewApp() (*App, error) {
	app := &App{
		Mode: createDefaultMode(),
	}

	gui, err := gui.NewGui()
	if err != nil {
		return nil, err
	}

	fm, err := fs.NewFileManager()
	if err != nil {
		return nil, err
	}

	app.Gui = gui
	app.FileManager = fm

	return app, nil
}

func (app *App) Run() error {
	go app.loop()

	return app.Gui.Run(app.onKey)
}

func (app *App) onModeChanged() {
	helps := make([]string, len(app.Mode.keyBindings.onKeys))
	idx := 0

	for k, a := range app.Mode.keyBindings.onKeys {
		helps[idx] = k + " " + a.help
		idx++
	}

	app.Gui.SetViewContent(app.Gui.Views.HelpMenu, helps)
}

func (app *App) onKey(key string) error {
	if action, hasKey := app.Mode.keyBindings.onKeys[key]; hasKey {
		for _, message := range action.messages {
			if err := message(app.Gui); err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *App) loop() {
	// Wait until Gui is loaded
	<-app.Gui.GuiLoadedChan
	// Load help menu
	app.onModeChanged()

	for {
		for range app.FileManager.DirLoadedChan {
			if err := app.Gui.Views.Main.SetCursor(0, 1); err != nil {
				log.Printf("failed to set cursor directory %v", err)
			}

			nodeSize := len(app.FileManager.Dir.Nodes)
			app.Gui.State.Main.SelectedIdx = 1
			app.Gui.State.Main.NumberOfFiles = nodeSize

			app.Gui.Views.Main.Title = " " + app.FileManager.Dir.Path +
				" (" + strconv.Itoa(nodeSize) + ") "
			lines := make([]string, nodeSize+1)

			lines[0] = "╭──── path"

			for i, node := range app.FileManager.Dir.Nodes {
				if i == nodeSize-1 {
					lines[i+1] = "╰─" + "  " + node.RelativePath
				} else {
					lines[i+1] = "├─" + "  " + node.RelativePath
				}
			}

			app.Gui.SetViewContent(app.Gui.Views.Main, lines)
		}
	}
}

func createDefaultMode() *Mode {
	return &Mode{
		name: "default",
		keyBindings: &KeyBindings{
			onKeys: map[string]*Action{
				"j": {
					help: "down",
					messages: []func(gui *gui.Gui) error{
						focusNext,
					},
				},
				"k": {
					help: "up",
					messages: []func(gui *gui.Gui) error{
						focusPrevious,
					},
				},
			},
		},
	}
}
