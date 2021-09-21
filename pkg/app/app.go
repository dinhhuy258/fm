package app

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/message"
	"github.com/dinhhuy258/fm/pkg/mode"
	"github.com/dinhhuy258/fm/pkg/state"
)

type App struct {
	gui         *gui.Gui
	fileManager *fs.FileManager
	state       *state.State
	modes       *mode.Modes
}

// NewApp bootstrap a new application
func NewApp() (*App, error) {
	app := &App{
		state: &state.State{
			FocusIdx:      0,
			NumberOfFiles: 0,
			Selections:    map[string]struct{}{},
		},
	}

	app.modes = mode.NewModes()

	g, err := gui.NewGui(app.onViewsCreated)
	if err != nil {
		return nil, err
	}

	app.gui = g

	return app, nil
}

func (app *App) Run() error {
	return app.gui.Run()
}

func (app *App) onModeChanged() {
	keys := make([]string, 0, len(app.modes.Peek().KeyBindings.OnKeys))
	helps := make([]string, 0, len(app.modes.Peek().KeyBindings.OnKeys))

	for k, a := range app.modes.Peek().KeyBindings.OnKeys {
		keys = append(keys, k)
		helps = append(helps, a.Help)
	}

	app.gui.Views.Help.SetTitle(app.modes.Peek().Name)

	if err := app.gui.Views.Help.SetHelp(keys, helps); err != nil {
		log.Fatalf("failed to set content for help view %v", err)
	}
}

func (app *App) State() *state.State {
	return app.state
}

func (app *App) Gui() *gui.Gui {
	return app.gui
}

func (app *App) FileManager() *fs.FileManager {
	return app.fileManager
}

func (app *App) PopMode() error {
	if err := app.modes.Pop(); err != nil {
		return err
	}

	app.onModeChanged()

	return nil
}

func (app *App) PushMode(mode string) error {
	if err := app.modes.Push(mode); err != nil {
		return err
	}

	app.onModeChanged()

	return nil
}

func (app *App) onKey(key string) error {
	if action, hasKey := app.modes.Peek().KeyBindings.OnKeys[key]; hasKey {
		for _, m := range action.Messages {
			if err := m.Func(app, m.Args...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *App) onViewsCreated() {
	// Load help menu
	if err := app.PushMode("default"); err != nil {
		log.Fatalf("failed to push default mode %v", err)
	}

	// Set on key handler
	app.gui.SetOnKeyFunc(app.onKey)

	fm, err := fs.NewFileManager()
	if err != nil {
		log.Fatalf("failed to create new file manager %v", err)
	}

	app.state.History = state.NewHistory(fm.Dir.Path)

	app.fileManager = fm

	go app.loop()
}

func (app *App) loop() {
	for {
		for range app.fileManager.DirLoadedChan {
			if err := app.gui.Views.Main.SetOrigin(0, 0); err != nil {
				log.Fatalf("failed to set origin %v", err)
			}

			if err := app.gui.Views.Main.SetCursor(0, 0); err != nil {
				log.Fatalf("failed to set cursor %v", err)
			}

			nodeSize := len(app.fileManager.Dir.VisibleNodes)
			app.state.FocusIdx = 0
			app.state.NumberOfFiles = nodeSize

			app.gui.Views.Main.SetTitle(" " + app.fileManager.Dir.Path + " (" + strconv.Itoa(nodeSize) + ") ")

			lastPath := app.state.History.Peek()
			if filepath.Dir(lastPath) == app.fileManager.Dir.Path {
				// back
				if err := message.FocusPath(app, lastPath); err != nil {
					log.Fatalf("failed to focus path %v", err)
				}
			}

			if err := app.gui.Views.Main.RenderDir(
				app.fileManager.Dir,
				app.state.Selections,
				app.state.FocusIdx,
			); err != nil {
				log.Fatalf("failed to render dir %v", err)
			}
		}
	}
}
