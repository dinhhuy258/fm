package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/key"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

// IApp is the interface for the application.
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

// Message is the message type.
type Message struct {
	Func func(app IApp, params ...string)
	Args []string
}
