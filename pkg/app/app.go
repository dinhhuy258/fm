package app

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/alitto/pond"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/fm/pkg/msg"
	"github.com/dinhhuy258/fm/pkg/pipe"
	"github.com/dinhhuy258/gocui"
)

var messageInRegexp = regexp.MustCompile(`[^\s']+|'([^']*)'`)

type App struct {
	gui        *gui.Gui
	modes      *Modes
	pressedKey key.Key
	pipe       *pipe.Pipe

	messageWorkerPool *pond.WorkerPool
}

// NewApp bootstrap a new application
func NewApp() (*App, error) {
	pipe, err := pipe.NewPipe()
	if err != nil {
		return nil, err
	}

	gui, err := gui.NewGui(pipe)
	if err != nil {
		return nil, err
	}

	app := &App{
		gui:               gui,
		pipe:              pipe,
		messageWorkerPool: pond.New(1 /* we only need one worker to avoid concurrency issue */, 10),
	}

	app.modes = CreateModes()

	return app, nil
}

func (app *App) Run() error {
	// Start watcher for the pipe
	app.pipe.StartWatcher(app.onMessageIn)
	// Push the default mode
	// TODO: Remove hard code here
	app.PushMode("default")

	// Set on key handler
	app.gui.SetOnKeyFunc(app.onKey)

	wd, err := os.Getwd()
	if err != nil {
		return nil
	}

	msg.ChangeDirectory(app, wd)

	return app.gui.Run()
}

func (app *App) GetPipe() *pipe.Pipe {
	return app.pipe
}

func (app *App) OnUIThread(f func() error) {
	app.gui.OnUIThread(f)
}

func (app *App) onModeChanged() {
	currentMode := app.modes.Peek()

	helps := currentMode.GetHelp()

	helpKeys := make([]string, 0, len(helps))
	helpMsgs := make([]string, 0, len(helps))

	for _, h := range helps {
		helpKeys = append(helpKeys, key.GetKeyDisplay(h.Key))
		helpMsgs = append(helpMsgs, h.Msg)
	}

	helpController, _ := app.GetController(controller.Help).(*controller.HelpController)
	helpController.SetHelp(currentMode.GetName(), helpKeys, helpMsgs)
	helpController.UpdateView()
}

func (app *App) GetController(controllerType controller.Type) controller.IController {
	return app.gui.GetController(controllerType)
}

func (app *App) PopMode() {
	if err := app.modes.Pop(); err != nil {
		// TODO: Better error handling???
		log.Fatalf("failed to pop mode %v", err)
	}

	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	logController.SetVisible(true)

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

func (app *App) Suspend() error {
	return app.gui.Suspend()
}

func (app *App) Resume() error {
	return app.gui.Resume()
}

func (app *App) OnQuit() {
	app.pipe.StopWatcher()
}

func (app *App) onKey(k gocui.Key, ch rune, _ gocui.Modifier) error {
	keybindings := app.modes.Peek().GetKeyBindings()

	if ch == 0 {
		app.pressedKey = k
	} else {
		app.pressedKey = ch
	}

	action, hasKey := keybindings.OnKeys[app.pressedKey]

	switch {
	case hasKey:
		app.submitMessages(action.Messages)
	case keybindings.Default != nil:
		app.submitMessages(keybindings.Default.Messages)
	}

	return nil
}

// submitMessages submit messages to the message worker pool
func (app *App) submitMessages(messages []*msg.Message) {
	for _, message := range messages {
		message := message // This will make scopelint happy

		app.messageWorkerPool.Submit(func() {
			message.Func(app, message.Args...)
		})
	}

	// Request re-render the GUI after each action
	app.messageWorkerPool.Submit(func() {
		app.gui.Render()
	})
}

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

	app.submitMessages([]*msg.Message{message})
}
