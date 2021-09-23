package message

import "github.com/dinhhuy258/fm/pkg/ctx"

func FocusNext(ctx ctx.Context, _ ...interface{}) error {
	if ctx.State().FocusIdx == ctx.State().NumberOfFiles-1 {
		return nil
	}

	if err := ctx.Gui().Views.Main.NextCursor(); err != nil {
		return err
	}

	ctx.State().FocusIdx++

	ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

func FocusPrevious(ctx ctx.Context, _ ...interface{}) error {
	if ctx.State().FocusIdx == 0 {
		return nil
	}

	if err := ctx.Gui().Views.Main.PreviousCursor(); err != nil {
		return err
	}

	ctx.State().FocusIdx--

	ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

func FocusPath(ctx ctx.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParameter
	}

	path, ok := params[0].(string)
	if !ok {
		return ErrInvalidMessageParameter
	}

	focusIdx := 0

	for _, node := range ctx.FileManager().Dir.VisibleNodes {
		if node.AbsolutePath == path {
			break
		}

		focusIdx++
	}

	if focusIdx == len(ctx.FileManager().Dir.VisibleNodes) {
		return nil
	}

	if err := ctx.Gui().Views.Main.SetCursor(0, 0); err != nil {
		return err
	}

	ctx.State().FocusIdx = 0

	for i := 0; i < focusIdx; i++ {
		if err := ctx.Gui().Views.Main.NextCursor(); err != nil {
			return err
		}

		ctx.State().FocusIdx++
	}

	ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

func Enter(ctx ctx.Context, _ ...interface{}) error {
	currentNode := ctx.FileManager().Dir.VisibleNodes[ctx.State().FocusIdx]

	if currentNode.IsDir {
		changeDirectory(ctx, currentNode.AbsolutePath, true)
	}

	return nil
}

func Back(ctx ctx.Context, _ ...interface{}) error {
	parent := ctx.FileManager().Dir.Parent()

	changeDirectory(ctx, parent, true)

	return nil
}

func LastVisitedPath(ctx ctx.Context, _ ...interface{}) error {
	ctx.State().History.VisitLast()
	changeDirectory(ctx, ctx.State().History.Peek(), false)

	return nil
}

func NextVisitedPath(ctx ctx.Context, _ ...interface{}) error {
	ctx.State().History.VisitNext()
	changeDirectory(ctx, ctx.State().History.Peek(), false)

	return nil
}

func changeDirectory(ctx ctx.Context, path string, saveHistory bool) {
	if saveHistory {
		ctx.State().History.Push(ctx.FileManager().Dir.Path)
	}

	ctx.FileManager().LoadDirectory(path)
}
