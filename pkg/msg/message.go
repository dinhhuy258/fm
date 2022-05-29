package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

type IApp interface {
	GetPipe() *pipe.Pipe
	// Controller
	GetController(controller.Type) controller.IController
	// Key
	GetPressedKey() key.Key
	// Quit
	Quit()
	// Mode
	PopMode()
	PushMode(mode string)
	// GUI
	OnUIThread(f func() error)
	Resume() error
	Suspend() error
}

type Message struct {
	Func func(app IApp, params ...string)
	Args []string
}
