package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
)

type IApp interface {
	// Render
	RenderEntries()
	RenderSelections()
	// Selection
	ClearSelections()
	GetSelections() []string
	ToggleSelection(path string)
	// Focus
	GetFocus() int
	SetFocus(focus int)
	// History
	PushHistory(entry fs.IEntry)
	PeekHistory() fs.IEntry
	VisitLastHistory()
	VisitNextHistory()
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
