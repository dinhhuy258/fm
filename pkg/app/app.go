package app

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/message"
	"github.com/dinhhuy258/fm/pkg/mode"
	"github.com/dinhhuy258/fm/pkg/state"
)

type App struct {
	Gui         *gui.Gui
	FileManager *fs.FileManager
	State       *state.State
	Modes       *mode.Modes
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

	app.Modes = mode.NewModes()

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
	keys := make([]string, 0, len(app.Modes.Peek().KeyBindings.OnKeys))
	helps := make([]string, 0, len(app.Modes.Peek().KeyBindings.OnKeys))

	for k, a := range app.Modes.Peek().KeyBindings.OnKeys {
		keys = append(keys, k)
		helps = append(helps, a.Help)
	}

	app.Gui.Views.Help.SetTitle(app.Modes.Peek().Name)

	if err := app.Gui.Views.Help.SetHelp(keys, helps); err != nil {
		log.Fatalf("failed to set content for help view %v", err)
	}
}

func (app *App) GetState() *state.State {
	return app.State
}

func (app *App) GetGui() *gui.Gui {
	return app.Gui
}

func (app *App) GetFileManager() *fs.FileManager {
	return app.FileManager
}

func (app *App) PopMode() error {
	if err := app.Modes.Pop(); err != nil {
		return err
	}

	app.onModeChanged()

	return nil
}

func (app *App) PushMode(mode string) error {
	if err := app.Modes.Push(mode); err != nil {
		return err
	}

	app.onModeChanged()

	return nil
}

func (app *App) onKey(key string) error {
	if action, hasKey := app.Modes.Peek().KeyBindings.OnKeys[key]; hasKey {
		var ctx ctx.Context = app
		for _, message := range action.Messages {
			if err := message.Func(&ctx, message.Args...); err != nil {
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
				var ctx ctx.Context = app
				if err := message.Focus(&ctx, lastPath); err != nil {
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
