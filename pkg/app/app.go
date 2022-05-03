package app

import (
	"log"
	"os"

	"github.com/dinhhuy258/fm/pkg/app/command"
	"github.com/dinhhuy258/fm/pkg/app/context"
	"github.com/dinhhuy258/fm/pkg/app/mode"
	"github.com/dinhhuy258/fm/pkg/fs"
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

	appGui := gui.GetGui()
	appGui.Views.Help.SetTitle(currentMode.GetName())
	appGui.Views.Help.SetHelp(keys, helps)
}

func (app *App) State() *context.State {
	return app.state
}

func (app *App) ClearSelections() {
	for k := range app.state.Selections {
		delete(app.state.Selections, k)
	}
}

func (app *App) DeleteSelection(path string) {
	delete(app.state.Selections, path)
}

func (app *App) HasSelection(path string) bool {
	_, hasSelection := app.state.Selections[path]

	return hasSelection
}

func (app *App) AddSelection(path string) {
	app.state.Selections[path] = struct{}{}
}

func (app *App) GetSelections() map[string]struct{} {
	return app.state.Selections
}

func (app *App) GetFocusIdx() int {
	return app.state.FocusIdx
}

func (app *App) SetFocusIdx(idx int) {
	app.state.FocusIdx = idx
}

func (app *App) GetNumberOfFiles() int {
	return app.state.NumberOfFiles
}

func (app *App) SetNumberOfFiles(numberOfFiles int) {
	app.state.NumberOfFiles = numberOfFiles
}

func (app *App) PushHistory(node *fs.Node) {
	app.State().History.Push(node)
}

func (app *App) PeekHistory() *fs.Node {
	return app.State().History.Peek()
}

func (app *App) VisitLastHistory() {
	app.State().History.VisitLast()
}

func (app *App) VisitNextHistory() {
	app.State().History.VisitNext()
}

func (app *App) MarkSave(key, path string) {
	app.State().Marks[key] = path
}

func (app *App) MarkLoad(key string) (string, bool) {
	path, hasKey := app.State().Marks[key]

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

	command.ChangeDirectory(app, wd, true, nil)
}
