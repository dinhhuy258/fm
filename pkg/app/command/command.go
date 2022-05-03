package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
)

type IApp interface {
	ClearSelections()
	GetSelections() map[string]struct{}
	DeleteSelection(path string)
	HasSelection(path string) bool
	AddSelection(path string)
	GetFocusIdx() int
	SetFocusIdx(idx int)
	GetNumberOfFiles() int
	SetNumberOfFiles(numberOfFiles int)
	PushHistory(node *fs.Node)
	PeekHistory() *fs.Node
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
