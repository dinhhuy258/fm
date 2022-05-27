package message

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/key"
)

type IApp interface {
	// Controller
	GetController(controller.Type) controller.IController
	// Key
	GetPressedKey() key.Key
	// Quit
	Quit()
	// Mark
	MarkSave(key, path string)
	MarkLoad(key string) (string, bool)
	// Mode
	PopMode()
	PushMode(mode string)
	// UI
	Resume()
	Suspend()
}

type Message struct {
	Func func(app IApp, params ...string)
	Args []string
}
