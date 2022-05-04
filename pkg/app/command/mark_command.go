package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
)

func MarkSave(app IApp, params ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()

	key, _ := params[0].(string)
	// Exit mark mode
	_ = app.PopMode()
	entry := fileExplorer.GetEntry(app.GetFocusIdx())
	app.MarkSave(key, entry.GetPath())

	return nil
}

func MarkLoad(app IApp, params ...interface{}) error {
	key, _ := params[0].(string)
	// Exit mark mode
	_ = app.PopMode()

	if path, hasKey := app.MarkLoad(key); hasKey {
		return FocusPath(app, path)
	}

	return nil
}
