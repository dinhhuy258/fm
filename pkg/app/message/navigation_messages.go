package message

import (
	"github.com/dinhhuy258/fm/pkg/app/context"
	"path/filepath"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

func FocusFirst(ctx context.Context, _ ...interface{}) error {
	_ = gui.GetGui().Views.Main.SetOrigin(0, 0)
	_ = gui.GetGui().Views.Main.SetCursor(0, 0)
	ctx.State().FocusIdx = 0

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

func FocusNext(ctx context.Context, _ ...interface{}) error {
	if ctx.State().FocusIdx == ctx.State().NumberOfFiles-1 {
		return nil
	}

	if err := gui.GetGui().Views.Main.NextCursor(); err != nil {
		return err
	}

	ctx.State().FocusIdx++

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

func FocusPrevious(ctx context.Context, _ ...interface{}) error {
	if ctx.State().FocusIdx == 0 {
		return nil
	}

	if err := gui.GetGui().Views.Main.PreviousCursor(); err != nil {
		return err
	}

	ctx.State().FocusIdx--

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

func FocusPath(ctx context.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParameter
	}

	path, ok := params[0].(string)
	if !ok {
		return ErrInvalidMessageParameter
	}

	if fs.GetFileManager().Dir.Path != filepath.Dir(path) {
		dirLoadedChan := fs.GetFileManager().LoadDirectory(filepath.Dir(path))
		<-dirLoadedChan
	}

	focusIdx := 0

	for _, node := range fs.GetFileManager().Dir.VisibleNodes {
		if node.AbsolutePath == path {
			break
		}

		focusIdx++
	}

	if focusIdx == len(fs.GetFileManager().Dir.VisibleNodes) {
		focusIdx = len(fs.GetFileManager().Dir.VisibleNodes) - 1
	}

	_ = gui.GetGui().Views.Main.SetCursor(0, 0)
	_ = gui.GetGui().Views.Main.SetOrigin(0, 0)

	ctx.State().FocusIdx = 0

	for i := 0; i < focusIdx; i++ {
		if err := gui.GetGui().Views.Main.NextCursor(); err != nil {
			return err
		}

		ctx.State().FocusIdx++
	}

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

func Enter(ctx context.Context, _ ...interface{}) error {
	currentNode := fs.GetFileManager().Dir.VisibleNodes[ctx.State().FocusIdx]

	if currentNode.IsDir {
		ChangeDirectory(ctx, currentNode.AbsolutePath, true, nil)
	}

	return nil
}

func Back(ctx context.Context, _ ...interface{}) error {
	parent := fs.GetFileManager().Dir.Parent()

	ChangeDirectory(ctx, parent, true, &fs.GetFileManager().Dir.Path)

	return nil
}

func LastVisitedPath(ctx context.Context, _ ...interface{}) error {
	ctx.State().History.VisitLast()
	node := ctx.State().History.Peek()
	ChangeDirectory(ctx, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func NextVisitedPath(ctx context.Context, _ ...interface{}) error {
	ctx.State().History.VisitNext()
	node := ctx.State().History.Peek()
	ChangeDirectory(ctx, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func ChangeDirectory(ctx context.Context, path string, saveHistory bool, focusPath *string) {
	if saveHistory && fs.GetFileManager().Dir != nil {
		currentNode := fs.GetFileManager().Dir.VisibleNodes[ctx.State().FocusIdx]
		ctx.State().History.Push(currentNode)
	}

	dirLoadedChan := fs.GetFileManager().LoadDirectory(path)

	go func() {
		<-dirLoadedChan

		numberOfFiles := len(fs.GetFileManager().Dir.VisibleNodes)
		ctx.State().NumberOfFiles = numberOfFiles
		gui.GetGui().Views.Main.SetTitle(fs.GetFileManager().Dir.Path, numberOfFiles)

		if focusPath == nil {
			_ = FocusFirst(ctx)
		} else {
			_ = FocusPath(ctx, *focusPath)
		}

		if saveHistory {
			currentNode := fs.GetFileManager().Dir.VisibleNodes[ctx.State().FocusIdx]
			ctx.State().History.Push(currentNode)
		}
	}()
}
