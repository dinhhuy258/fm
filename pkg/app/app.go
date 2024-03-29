package app

import (
	"os"
	"regexp"
	"strings"

	"github.com/alitto/pond"
	"github.com/dinhhuy258/gocui"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/config/lua"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/key"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/msg"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

var messageInRegexp = regexp.MustCompile(`[^\s']+|'([^']*)'`)

const (
	// The number of worker in the message worker pool should be 1 to avoid concurrency issue
	maxMessageWorkers = 1
	// The maximum number of messages that can be queued
	maxMessageCapacity = 10
)

// App is the main application
type App struct {
	gui   *gui.Gui
	lua   *lua.Lua
	modes *Modes
	pipe  *pipe.Pipe

	messageWorkerPool *pond.WorkerPool
}

// NewApp bootstrap a new application
func NewApp() (*App, error) {
	lua := lua.NewLua()

	// Load the config
	if err := config.LoadConfig(lua); err != nil {
		return nil, err
	}

	pipe, err := pipe.NewPipe()
	if err != nil {
		return nil, err
	}

	gui, err := gui.NewGui()
	if err != nil {
		return nil, err
	}

	app := &App{
		gui:               gui,
		lua:               lua,
		pipe:              pipe,
		messageWorkerPool: pond.New(maxMessageWorkers, maxMessageCapacity),
	}

	app.modes = CreateModes(app.onModeChange)

	return app, nil
}

// Run the application
func (app *App) Run() error {
	// Start watcher for the pipe
	app.pipe.StartWatcher(app.onMessageIn)

	// Set on key handler
	app.gui.SetOnKeyFunc(app.onKey)

	// Get the current directory and load files/folder in it
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	ctx := make(msg.MessageContext)
	ctx["arg1"] = wd
	msg.ChangeDirectory(app, nil, ctx)

	return app.gui.Run()
}

// GetPipe returns the pipe
func (app *App) GetPipe() *pipe.Pipe {
	return app.pipe
}

// OnUIThread is called to handle messages from the UI thread
func (app *App) OnUIThread(f func() error) {
	app.gui.OnUIThread(f)
}

// onModeChange is called when the mode is changed
func (app *App) onModeChange(currentMode *Mode) {
	helps := currentMode.GetHelp()

	helpKeys := make([]string, 0, len(helps))
	helpMsgs := make([]string, 0, len(helps))

	for _, h := range helps {
		helpKeys = append(helpKeys, key.GetKeyDisplay(h.key))
		helpMsgs = append(helpMsgs, h.msg)
	}

	helpController, _ := app.GetController(controller.Help).(*controller.HelpController)
	helpController.SetHelp(currentMode.GetName(), helpKeys, helpMsgs)
	helpController.UpdateView()
}

// GetController returns the controller with the given name
func (app *App) GetController(controllerType controller.Type) controller.IController {
	return app.gui.GetController(controllerType)
}

func (app *App) SetLog(level view.LogLevel, msgFormat string, args ...interface{}) {
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	logController.SetLog(level, msgFormat, args...)
	logController.UpdateView()
}

// SwitchMode switches to the given mode
func (app *App) SwitchMode(mode string) {
	if err := app.modes.SwitchMode(mode); err != nil {
		app.SetLog(view.LogError, "Mode not found: "+mode)
	}

	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	// make sure that the log view is always show after switching mode
	logController.ShowLog()
}

// Quit the application
func (app *App) Quit() {
	app.gui.Quit()
}

// Suspend the application
func (app *App) Suspend() error {
	return app.gui.Suspend()
}

// Resume the application
func (app *App) Resume() error {
	return app.gui.Resume()
}

// OnQuit is called when the application is about to quit
func (app *App) OnQuit() {
	app.lua.Close()
	app.pipe.StopWatcher()
}

// onKey is called from gocui when a key is pressed
func (app *App) onKey(k gocui.Key, ch rune, mod gocui.Modifier) error {
	keybindings := app.modes.GetCurrentMode().GetKeyBindings()

	pressedKey := key.Key{
		Key: k,
		Ch:  ch,
		Mod: mod,
	}

	action, hasKey := keybindings.onKeys[pressedKey]

	switch {
	case hasKey:
		app.submitMessages(action.messages, &pressedKey)
	case pressedKey.IsDigit() && keybindings.onNumber != nil:
		app.submitMessages(keybindings.onNumber.messages, &pressedKey)
	case keybindings.defaultAction != nil:
		app.submitMessages(keybindings.defaultAction.messages, &pressedKey)
	default:
		app.SetLog(view.LogWarning, "Key map not found")
	}

	return nil
}

// submitMessages submit messages to the message worker pool
func (app *App) submitMessages(messages []*msg.Message, pressedKey *key.Key) {
	for _, message := range messages {
		message := message // This will make scopelint happy

		app.messageWorkerPool.Submit(func() {
			message.Func(app, pressedKey, message.Ctx)
		})
	}

	// Request re-render the GUI after each action
	app.messageWorkerPool.Submit(func() {
		app.gui.Render()
	})
}

// onMessageIn is called when a message is received from the pipe
func (app *App) onMessageIn(messageIn string) {
	if messageIn == "" {
		// Ignore empty message
		return
	}

	components := messageInRegexp.FindAllString(messageIn, -1)
	if len(components) == 0 {
		return
	}

	args := components[1:]
	for idx, arg := range args {
		args[idx] = strings.Trim(arg, "'")
	}

	message, err := msg.NewMessage(components[0], args...)
	if err != nil {
		return
	}

	app.submitMessages([]*msg.Message{message}, nil)
}
