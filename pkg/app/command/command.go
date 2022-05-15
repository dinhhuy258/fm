package command

import "github.com/dinhhuy258/fm/pkg/gui"

type IApp interface {
	// Gui
	GetGui() *gui.Gui
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
