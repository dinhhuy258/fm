package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
)

type IApp interface {
	RenderEntries()
	RenderSelections()
	ClearSelections()
	GetSelections() map[string]struct{}
	DeleteSelection(path string)
	HasSelection(path string) bool
	AddSelection(path string)
	GetFocusIdx() int
	SetFocusIdx(idx int)
	PushHistory(entry fs.IEntry)
	PeekHistory() fs.IEntry
	VisitLastHistory()
	VisitNextHistory()
	MarkSave(key, path string)
	MarkLoad(key string) (string, bool)
	PopMode() error
	PushMode(mode string) error
}

type Command struct {
	Help string
	Func func(app IApp, params ...interface{}) error
	Args []interface{}
}
