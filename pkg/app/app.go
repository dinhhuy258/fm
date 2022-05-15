package app

import (
	"log"
	"os"

	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/gui"
)

type App struct {
	gui   *gui.Gui
	Marks map[string]string
	modes *Modes
}

// NewApp bootstrap a new application
func NewApp() *App {
	app := &App{
		gui:   gui.NewGui(),
		Marks: map[string]string{},
	}

	app.modes = NewModes()

	return app
}

func (app *App) Run() error {
	return app.gui.Run(app.onGuiReady)
}

func (app *App) onModeChanged() {
	currentMode := app.modes.Peek()

	currentMode.OnModeStarted(app)

	helps := currentMode.GetHelp(app)

	keys := make([]string, 0, len(helps))
	msgs := make([]string, 0, len(helps))

	for _, h := range helps {
		keys = append(keys, h.Key)
		msgs = append(msgs, h.Msg)
	}

	app.gui.GetControllers().Help.SetHelp(currentMode.GetName(), keys, msgs)
}

func (app *App) GetGui() *gui.Gui {
	return app.gui
}

func (app *App) MarkSave(key, path string) {
	app.Marks[key] = path
}

func (app *App) MarkLoad(key string) (string, bool) {
	path, hasKey := app.Marks[key]

	return path, hasKey
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
	keybindings := app.modes.Peek().GetKeyBindings()

	if action, hasKey := keybindings.OnKeys[key]; hasKey {
		for _, cmd := range action.Commands {
			if err := cmd.Func(app, cmd.Args...); err != nil {
				return err
			}
		}
	} else if keybindings.OnAlphabet != nil {
		for _, cmd := range keybindings.OnAlphabet.Commands {
			args := cmd.Args
			args = append(args, key)

			if err := cmd.Func(app, args...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *App) onGuiReady() {
	// Push the default mode
	_ = app.PushMode("default")

	// Set on key handler
	app.gui.SetOnKeyFunc(app.onKey)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory %v", err)
	}

	command.ChangeDirectory(app, wd, "")
}
