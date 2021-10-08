package message

import (
	"path/filepath"

	"github.com/dinhhuy258/fm/pkg/ctx"
)

func FocusFirst(ctx ctx.Context, _ ...interface{}) error {
	_ = ctx.Gui().Views.Main.SetOrigin(0, 0)
	_ = ctx.Gui().Views.Main.SetCursor(0, 0)
	ctx.State().FocusIdx = 0

	ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

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

	if ctx.FileManager().Dir.Path != filepath.Dir(path) {
		dirLoadedChan := ctx.FileManager().LoadDirectory(filepath.Dir(path))
		<-dirLoadedChan
	}

	focusIdx := 0

	for _, node := range ctx.FileManager().Dir.VisibleNodes {
		if node.AbsolutePath == path {
			break
		}

		focusIdx++
	}

	if focusIdx == len(ctx.FileManager().Dir.VisibleNodes) {
		focusIdx = len(ctx.FileManager().Dir.VisibleNodes) - 1
	}

	_ = ctx.Gui().Views.Main.SetCursor(0, 0)
	_ = ctx.Gui().Views.Main.SetOrigin(0, 0)

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
		ChangeDirectory(ctx, currentNode.AbsolutePath, true, nil)
	}

	return nil
}

func Back(ctx ctx.Context, _ ...interface{}) error {
	parent := ctx.FileManager().Dir.Parent()

	ChangeDirectory(ctx, parent, true, &ctx.FileManager().Dir.Path)

	return nil
}

func LastVisitedPath(ctx ctx.Context, _ ...interface{}) error {
	ctx.State().History.VisitLast()
	node := ctx.State().History.Peek()
	ChangeDirectory(ctx, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func NextVisitedPath(ctx ctx.Context, _ ...interface{}) error {
	ctx.State().History.VisitNext()
	node := ctx.State().History.Peek()
	ChangeDirectory(ctx, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func ChangeDirectory(ctx ctx.Context, path string, saveHistory bool, focusPath *string) {
	if saveHistory && ctx.FileManager().Dir != nil {
		currentNode := ctx.FileManager().Dir.VisibleNodes[ctx.State().FocusIdx]
		ctx.State().History.Push(currentNode)
	}

	dirLoadedChan := ctx.FileManager().LoadDirectory(path)

	go func() {
		<-dirLoadedChan

		numberOfFiles := len(ctx.FileManager().Dir.VisibleNodes)
		ctx.State().NumberOfFiles = numberOfFiles
		ctx.Gui().Views.Main.SetTitle(ctx.FileManager().Dir.Path, numberOfFiles)

		if focusPath == nil {
			_ = FocusFirst(ctx)
		} else {
			_ = FocusPath(ctx, *focusPath)
		}

		if saveHistory {
			currentNode := ctx.FileManager().Dir.VisibleNodes[ctx.State().FocusIdx]
			ctx.State().History.Push(currentNode)
		}
	}()
}
