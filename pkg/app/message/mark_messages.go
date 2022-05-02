package message

import (
	"github.com/dinhhuy258/fm/pkg/app/context"
	"github.com/dinhhuy258/fm/pkg/fs"
)

func MarkSave(ctx context.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParameter
	}

	key, ok := params[0].(string)
	if !ok {
		return ErrInvalidMessageParameter
	}

	_ = ctx.PopMode()
	currentNode := fs.GetFileManager().Dir.VisibleNodes[ctx.State().FocusIdx]
	ctx.State().Marks[key] = currentNode.AbsolutePath

	return nil
}

func MarkLoad(ctx context.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParameter
	}

	key, ok := params[0].(string)
	if !ok {
		return ErrInvalidMessageParameter
	}

	_ = ctx.PopMode()

	if path, hasKey := ctx.State().Marks[key]; hasKey {
		return FocusPath(ctx, path)
	}

	return nil
}
