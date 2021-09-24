package message

import "github.com/dinhhuy258/fm/pkg/ctx"

func MarkSave(ctx ctx.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParameter
	}

	key, ok := params[0].(string)
	if !ok {
		return ErrInvalidMessageParameter
	}

	_ = ctx.PopMode()
	currentNode := ctx.FileManager().Dir.VisibleNodes[ctx.State().FocusIdx]
	ctx.State().Marks[key] = currentNode.RelativePath

	return nil
}
