package command

import (
	"github.com/dinhhuy258/fm/pkg/fs"
)

func MarkSave(app IApp, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidCommandParameter
	}

	key, ok := params[0].(string)
	if !ok {
		return ErrInvalidCommandParameter
	}

	_ = app.PopMode()
	currentNode := fs.GetFileManager().Dir.VisibleNodes[app.State().FocusIdx]
	app.State().Marks[key] = currentNode.AbsolutePath

	return nil
}

func MarkLoad(app IApp, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidCommandParameter
	}

	key, ok := params[0].(string)
	if !ok {
		return ErrInvalidCommandParameter
	}

	_ = app.PopMode()

	if path, hasKey := app.State().Marks[key]; hasKey {
		return FocusPath(app, path)
	}

	return nil
}
