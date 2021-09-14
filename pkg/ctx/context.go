package ctx

import (
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/state"
)

type Context interface {
	GetState() *state.State
	GetGui() *gui.Gui
	GetFileManager() *fs.FileManager
}
