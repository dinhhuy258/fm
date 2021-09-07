package app

import (
	"log"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/mode"
)

type App struct {
	Gui         *gui.Gui
	Mode        *mode.Mode
	FileManager *fs.FileManager
}

// NewApp bootstrap a new application
func NewApp() (*App, error) {
	app := &App{}

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

func (app *App) onKey(key string) error {
	app.Gui.SetViewContent(app.Gui.Views.Main, []string{key})

	return nil
}

func (app *App) loop() {
	// Wait until Gui is loaded
	<-app.Gui.GuiLoadedChan

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
