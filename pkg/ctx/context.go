package ctx

import (
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/state"
)

type Context interface {
	State() *state.State
	Gui() *gui.Gui
	FileManager() *fs.FileManager
	PopMode() error
	PushMode(mode string) error
}
