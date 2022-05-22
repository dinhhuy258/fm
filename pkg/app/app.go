package app

import (
	"log"
	"os"

	"github.com/alitto/pond"
	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/gocui"
)

type App struct {
	gui        *gui.Gui
	marks      map[string]string
	modes      *Modes
	pressedKey key.Key

	commandWorkerPool *pond.WorkerPool
}

// NewApp bootstrap a new application
func NewApp() *App {
	app := &App{
		gui:               gui.NewGui(),
		marks:             map[string]string{},
		commandWorkerPool: pond.New(1 /* we only need one woker to avoid concurency issue */, 10),
	}

	app.modes = CreateAllModes(app.marks)

	return app
}

func (app *App) Run() error {
	return app.gui.Run(app.onGuiReady)
}

func (app *App) onModeChanged() {
	currentMode := app.modes.Peek()

	helps := currentMode.GetHelp()

	keys := make([]string, 0, len(helps))
	msgs := make([]string, 0, len(helps))

	for _, h := range helps {
		keys = append(keys, key.GetKeyDisplay(h.Key))
		msgs = append(msgs, h.Msg)
	}

	helpController, _ := app.GetController(controller.Help).(*controller.HelpController)
	helpController.SetHelp(currentMode.GetName(), keys, msgs)
	helpController.UpdateView()
}

func (app *App) GetController(controllerType controller.Type) controller.IController {
	return app.gui.GetController(controllerType)
}

func (app *App) MarkSave(key, path string) {
	app.marks[key] = path
}

func (app *App) MarkLoad(key string) (string, bool) {
	path, hasKey := app.marks[key]

	return path, hasKey
}

func (app *App) PopMode() {
	if err := app.modes.Pop(); err != nil {
		// TODO: Better error handling???
		log.Fatalf("failed to pop mode %v", err)
	}

	app.onModeChanged()
}

func (app *App) PushMode(mode string) {
	if err := app.modes.Push(mode); err != nil {
		// TODO: Better error handling???
		log.Fatalf("failed to push mode %v", err)
	}

	app.onModeChanged()
}

func (app *App) GetPressedKey() key.Key {
	return app.pressedKey
}

func (app *App) Quit() {
	app.gui.Quit()
}

func (app *App) onKey(k gocui.Key, ch rune, _ gocui.Modifier) error {
	keybindings := app.modes.Peek().GetKeyBindings()

	if ch == 0 {
		app.pressedKey = k
	} else {
		app.pressedKey = ch
	}

	if action, hasKey := keybindings.OnKeys[app.pressedKey]; hasKey {
		for _, cmd := range action.Commands {
			cmd := cmd

			app.commandWorkerPool.Submit(func() {
				cmd.Func(app, cmd.Args...)
			})
		}
	} else if keybindings.OnAlphabet != nil {
		for _, cmd := range keybindings.OnAlphabet.Commands {
			cmd := cmd

			app.commandWorkerPool.Submit(func() {
				cmd.Func(app, cmd.Args...)
			})
		}
	} else if keybindings.Default != nil {
		for _, cmd := range keybindings.Default.Commands {
			cmd := cmd

			app.commandWorkerPool.Submit(func() {
				cmd.Func(app, cmd.Args...)
			})
		}
	}

	return nil
}

func (app *App) onGuiReady() {
	// Push the default mode
	app.PushMode("default")

	// Set on key handler
	app.gui.SetOnKeyFunc(app.onKey)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory %v", err)
	}

	command.ChangeDirectory(app, wd)
}
