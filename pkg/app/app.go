package app

import (
	"log"
	"os"

	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

type App struct {
	FocusIdx   int
	Selections map[string]struct{}
	History    *History
	Marks      map[string]string
	modes      *Modes
}

// NewApp bootstrap a new application
func NewApp() *App {
	app := &App{
		FocusIdx:   -1,
		Selections: map[string]struct{}{},
		Marks:      map[string]string{},
		History:    NewHistory(),
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
	keys, helps := currentMode.GetHelp(app)

	appGui := gui.GetGui()
	appGui.SetHelpTitle(currentMode.GetName())
	appGui.SetHelp(keys, helps)
}

func (app *App) ClearSelections() {
	for k := range app.Selections {
		delete(app.Selections, k)
	}
}

func (app *App) DeleteSelection(path string) {
	delete(app.Selections, path)
}

func (app *App) HasSelection(path string) bool {
	_, hasSelection := app.Selections[path]

	return hasSelection
}

func (app *App) AddSelection(path string) {
	app.Selections[path] = struct{}{}
}

func (app *App) GetSelections() map[string]struct{} {
	return app.Selections
}

func (app *App) GetFocusIdx() int {
	return app.FocusIdx
}

func (app *App) SetFocusIdx(idx int) {
	app.FocusIdx = idx
}

func (app *App) PushHistory(entry fs.IEntry) {
	app.History.Push(entry)
}

func (app *App) PeekHistory() fs.IEntry {
	return app.History.Peek()
}

func (app *App) VisitLastHistory() {
	app.History.VisitLast()
}

func (app *App) VisitNextHistory() {
	app.History.VisitNext()
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

	if cmd, hasKey := keybindings.OnKeys[key]; hasKey {
		if err := cmd.Func(app, cmd.Args...); err != nil {
			return err
		}
	} else if keybindings.OnAlphabet != nil {
		args := keybindings.OnAlphabet.Args
		args = append(args, key)

		if err := keybindings.OnAlphabet.Func(app, args...); err != nil {
			return err
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

	command.LoadDirectory(app, wd, true, "")
}
