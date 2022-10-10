package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/key"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

// IApp is the interface for the application.
type IApp interface {
	GetPipe() *pipe.Pipe
	// Controller
	GetController(controller.Type) controller.IController
	// Quit
	Quit()
	// Mode
	SwitchMode(mode string)
	// GUI
	OnUIThread(f func() error)
	Resume() error
	Suspend() error
	SetLog(level view.LogLevel, msgFormat string, args ...interface{})
}

// Message is the message type.
type Message struct {
	Func func(app IApp, key key.Key, ctx MessageContext)
	Ctx  MessageContext
}
