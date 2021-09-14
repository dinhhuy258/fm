package message

import "github.com/dinhhuy258/fm/pkg/ctx"

func FocusNext(ctx *ctx.Context, params ...interface{}) error {
	if (*ctx).GetState().FocusIdx == (*ctx).GetState().NumberOfFiles-1 {
		return nil
	}

	if err := (*ctx).GetGui().Views.Main.NextCursor(); err != nil {
		return err
	}

	(*ctx).GetState().FocusIdx++

	return (*ctx).GetGui().Views.Main.RenderDir(
		(*ctx).GetFileManager().Dir,
		(*ctx).GetState().Selections,
		(*ctx).GetState().FocusIdx,
	)
}

func FocusPrevious(ctx *ctx.Context, params ...interface{}) error {
	if (*ctx).GetState().FocusIdx == 0 {
		return nil
	}

	if err := (*ctx).GetGui().Views.Main.PreviousCursor(); err != nil {
		return err
	}

	(*ctx).GetState().FocusIdx--

	return (*ctx).GetGui().Views.Main.RenderDir(
		(*ctx).GetFileManager().Dir,
		(*ctx).GetState().Selections,
		(*ctx).GetState().FocusIdx,
	)
}

func Enter(ctx *ctx.Context, params ...interface{}) error {
	currentNode := (*ctx).GetFileManager().Dir.Nodes[(*ctx).GetState().FocusIdx]

	if currentNode.IsDir {
		changeDirectory(ctx, currentNode.AbsolutePath, true)
	}

	return nil
}

func Back(ctx *ctx.Context, params ...interface{}) error {
	parent := (*ctx).GetFileManager().Dir.Parent()

	changeDirectory(ctx, parent, true)

	return nil
}

func LastVisitedPath(ctx *ctx.Context, params ...interface{}) error {
	(*ctx).GetState().History.VisitLast()
	changeDirectory(ctx, (*ctx).GetState().History.Peek(), false)

	return nil
}

func NextVisitedPath(ctx *ctx.Context, params ...interface{}) error {
	(*ctx).GetState().History.VisitNext()
	changeDirectory(ctx, (*ctx).GetState().History.Peek(), false)

	return nil
}

func changeDirectory(ctx *ctx.Context, path string, saveHistory bool) {
	if saveHistory {
		(*ctx).GetState().History.Push((*ctx).GetFileManager().Dir.Path)
	}

	(*ctx).GetFileManager().LoadDirectory(path)
}
