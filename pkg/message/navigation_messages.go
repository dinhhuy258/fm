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

	return ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)
}

func FocusPrevious(ctx ctx.Context, _ ...interface{}) error {
	if ctx.State().FocusIdx == 0 {
		return nil
	}

	if err := ctx.Gui().Views.Main.PreviousCursor(); err != nil {
		return err
	}

	ctx.State().FocusIdx--

	return ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)
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
