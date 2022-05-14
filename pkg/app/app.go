package app

import (
	"log"
	"os"

	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/gui"
)

type App struct {
	Focus      int
	Selections map[string]struct{}
	Marks      map[string]string
	modes      *Modes
}

// NewApp bootstrap a new application
func NewApp() *App {
	app := &App{
		Focus:      -1,
		Selections: map[string]struct{}{},
		Marks:      map[string]string{},
	}

	app.modes = NewModes()

	gui.InitGui(app.onViewsCreated)

	return app
}

func (app *App) Run() error {
	return gui.GetGui().Run()
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

	appGui := gui.GetGui()
	appGui.GetControllers().Help.SetHelp(currentMode.GetName(), keys, msgs)
}

func (app *App) RenderEntries() {
	// fileExplorer := fs.GetFileExplorer()
	// appGui := gui.GetGui()

	// appGui.RenderEntries(
	// 	fileExplorer.GetEntries(),
	// 	app.Selections,
	// 	app.Focus,
	// )
}

func (app *App) RenderSelections() {
	appGui := gui.GetGui()

	appGui.RenderSelections(app.GetSelections())
}

func (app *App) ClearSelections() {
	for k := range app.Selections {
		delete(app.Selections, k)
	}
}

func (app *App) ToggleSelection(path string) {
	if _, hasSelection := app.Selections[path]; hasSelection {
		delete(app.Selections, path)
	} else {
		app.Selections[path] = struct{}{}
	}
}

func (app *App) GetSelections() []string {
	selections := make([]string, 0, len(app.Selections))
	for selection := range app.Selections {
		selections = append(selections, selection)
	}

	return selections
}

func (app *App) GetFocus() int {
	return app.Focus
}

func (app *App) SetFocus(focus int) {
	app.Focus = focus
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

func (app *App) onViewsCreated() {
	// Push the default mode
	_ = app.PushMode("default")

	// Set on key handler
	gui.GetGui().SetOnKeyFunc(app.onKey)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory %v", err)
	}

	command.LoadDirectory(app, wd, "")
}
