package app

import "github.com/dinhhuy258/fm/pkg/gui"

type App struct {
	Gui *gui.Gui
}

// NewApp bootstrap a new application
func NewApp() (*App, error) {
	app := &App{}

	gui, err := gui.NewGui()
	if err != nil {
		return nil, err
	}

	app.Gui = gui

	return app, nil
}

func (app App) Run() error {
	return app.Gui.Run()
}
