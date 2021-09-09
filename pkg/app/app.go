package app

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

type App struct {
	Gui         *gui.Gui
	History     *History
	State       *State
	Mode        *Mode
	FileManager *fs.FileManager
}

// NewApp bootstrap a new application
func NewApp() (*App, error) {
	app := &App{
		Mode: createDefaultMode(),
		State: &State{
			Main: &MainState{
				FocusIdx:      0,
				NumberOfFiles: 0,
			},
		},
	}

	gui, err := gui.NewGui()
	if err != nil {
		return nil, err
	}

	fm, err := fs.NewFileManager()
	if err != nil {
		return nil, err
	}

	app.History = NewHistory(fm.Dir.Path)

	app.Gui = gui
	app.FileManager = fm

	return app, nil
}

func (app *App) Run() error {
	go app.loop()

	return app.Gui.Run()
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
			if err := message(app); err != nil {
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
			app.State.Main.FocusIdx = 0
			app.State.Main.NumberOfFiles = nodeSize

			app.Gui.Views.MainHeader.Title = " " + app.FileManager.Dir.Path +
				" (" + strconv.Itoa(nodeSize) + ") "

			lastPath := app.History.Peek()
			if filepath.Dir(lastPath) == app.FileManager.Dir.Path {
				// back
				base := filepath.Base(lastPath)
				for _, node := range app.FileManager.Dir.Nodes {
					if node.IsDir && node.RelativePath == base {
						break
					}

					if err := app.Gui.NextCursor(app.Gui.Views.Main); err != nil {
						log.Fatalf("failed to focus next %v", err)
					}
					app.State.Main.FocusIdx++
				}
			}

			app.Gui.RenderDir(app.FileManager.Dir, app.State.Main.FocusIdx)
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
					messages: []func(app *App) error{
						focusNext,
					},
				},
				"k": {
					help: "up",
					messages: []func(app *App) error{
						focusPrevious,
					},
				},
				"l": {
					help: "enter",
					messages: []func(app *App) error{
						enter,
					},
				},
				"h": {
					help: "back",
					messages: []func(app *App) error{
						back,
					},
				},
				"ctrl+i": {
					help: "next visited path",
					messages: []func(app *App) error{
						nextVisitedPath,
					},
				},
				"ctrl+o": {
					help: "last visited path",
					messages: []func(app *App) error{
						lastVisitedPath,
					},
				},
				"q": {
					help: "quit",
					messages: []func(app *App) error{
						quit,
					},
				},
			},
		},
	}
}
