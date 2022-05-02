package app

import (
	"github.com/dinhhuy258/fm/pkg/app/context"
	"github.com/dinhhuy258/fm/pkg/app/message"
	"github.com/dinhhuy258/fm/pkg/app/mode"
	"log"
	"os"

	"github.com/dinhhuy258/fm/pkg/gui"
)

type App struct {
	state *context.State
	modes *mode.Modes
}

// NewApp bootstrap a new application
func NewApp() *App {
	app := &App{
		state: context.NewState(),
	}

	app.modes = mode.NewModes()

	gui.InitGui(app.onViewsCreated)

	return app
}

func (app *App) Run() error {
	return gui.GetGui().Run()
}

func (app *App) onModeChanged() {
	currentMode := app.modes.Peek()
	keys, helps := currentMode.GetHelp(app.state)

	gui := gui.GetGui()
	gui.Views.Help.SetTitle(currentMode.Name)
	gui.Views.Help.SetHelp(keys, helps)
}

func (app *App) State() *context.State {
	return app.state
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
	keybindings := app.modes.Peek().KeyBindings

	if action, hasKey := keybindings.OnKeys[key]; hasKey {
		for _, m := range action.Messages {
			if err := m.Func(app, m.Args...); err != nil {
				return err
			}
		}
	} else if keybindings.OnAlphabet != nil {
		for _, m := range keybindings.OnAlphabet.Messages {
			args := m.Args
			args = append(args, key)

			if err := m.Func(app, args...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *App) onViewsCreated() {
	// Load help menu
	_ = app.PushMode("default")

	// Set on key handler
	gui.GetGui().SetOnKeyFunc(app.onKey)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory %v", err)
	}

	message.ChangeDirectory(app, wd, true, nil)
}
