package app

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/state"
)

type App struct {
	Gui         *gui.Gui
	FileManager *fs.FileManager
	State       *state.State
	Modes       *Modes
}

// NewApp bootstrap a new application
func NewApp() (*App, error) {
	app := &App{
		State: &state.State{
			FocusIdx:      0,
			NumberOfFiles: 0,
			Selections:    map[string]struct{}{},
		},
	}

	app.Modes = NewModes()

	if err := app.Modes.Push("default"); err != nil {
		return nil, err
	}

	gui, err := gui.NewGui()
	if err != nil {
		return nil, err
	}

	fm, err := fs.NewFileManager()
	if err != nil {
		return nil, err
	}

	app.State.History = state.NewHistory(fm.Dir.Path)

	app.Gui = gui
	app.FileManager = fm

	return app, nil
}

func (app *App) Run() error {
	go app.loop()

	return app.Gui.Run()
}

func (app *App) onModeChanged() {
	keys := make([]string, 0, len(app.Modes.Peek().keyBindings.onKeys))
	helps := make([]string, 0, len(app.Modes.Peek().keyBindings.onKeys))

	for k, a := range app.Modes.Peek().keyBindings.onKeys {
		keys = append(keys, k)
		helps = append(helps, a.help)
	}

	app.Gui.Views.Help.SetTitle(app.Modes.Peek().name)

	if err := app.Gui.Views.Help.SetHelp(keys, helps); err != nil {
		log.Fatalf("failed to set content for help view %v", err)
	}
}

func (app *App) onKey(key string) error {
	if action, hasKey := app.Modes.Peek().keyBindings.onKeys[key]; hasKey {
		for _, message := range action.messages {
			if err := message.f(app, message.args...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *App) loop() {
	// Wait until Gui is loaded
	<-app.Gui.ViewsCreatedChan
	// Load help menu
	app.onModeChanged()
	// Set on key handler
	app.Gui.SetOnKeyFunc(app.onKey)

	for {
		for range app.FileManager.DirLoadedChan {
			if err := app.Gui.Views.Main.SetOrigin(0, 0); err != nil {
				log.Fatalf("failed to set origin %v", err)
			}

			if err := app.Gui.Views.Main.SetCursor(0, 0); err != nil {
				log.Fatalf("failed to set cursor %v", err)
			}

			nodeSize := len(app.FileManager.Dir.Nodes)
			app.State.FocusIdx = 0
			app.State.NumberOfFiles = nodeSize

			app.Gui.Views.Main.SetTitle(" " + app.FileManager.Dir.Path + " (" + strconv.Itoa(nodeSize) + ") ")

			lastPath := app.State.History.Peek()
			if filepath.Dir(lastPath) == app.FileManager.Dir.Path {
				// back
				if err := focus(app, lastPath); err != nil {
					log.Fatalf("failed to focus path %v", err)
				}
			}

			if err := app.Gui.Views.Main.RenderDir(
				app.FileManager.Dir,
				app.State.Selections,
				app.State.FocusIdx,
			); err != nil {
				log.Fatalf("failed to render dir %v", err)
			}
		}
	}
}
