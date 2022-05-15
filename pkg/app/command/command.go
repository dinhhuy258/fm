package command

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
)

type IApp interface {
	// Controller
	GetController(controller.Type) controller.IController
	// Quit
	Quit() error
	// Mark
	MarkSave(key, path string)
	MarkLoad(key string) (string, bool)
	// Mode
	PopMode() error
	PushMode(mode string) error
}

type Command struct {
	Func func(app IApp, params ...interface{}) error
	Args []interface{}
}
