package message

import (
	"errors"
	"github.com/dinhhuy258/fm/pkg/app/context"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/gocui"
)

var ErrInvalidMessageParameter = errors.New("invalid message parameter")

func ToggleSelection(ctx context.Context, _ ...interface{}) error {
	path := fs.GetFileManager().Dir.VisibleNodes[ctx.State().FocusIdx].AbsolutePath

	if _, hasPath := ctx.State().Selections[path]; hasPath {
		delete(ctx.State().Selections, path)
	} else {
		ctx.State().Selections[path] = struct{}{}
	}

	refreshSelections(ctx)

	return nil
}

func ToggleHidden(ctx context.Context, _ ...interface{}) error {
	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	fs.GetFileManager().Dir.Reload()

	numberOfFiles := len(fs.GetFileManager().Dir.VisibleNodes)
	ctx.State().NumberOfFiles = numberOfFiles
	gui.GetGui().Views.Main.SetTitle(fs.GetFileManager().Dir.Path, numberOfFiles)
	gui.GetGui().Views.SortAndFilter.SetSortAndFilter()

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

func ClearSelection(ctx context.Context, _ ...interface{}) error {
	ctx.State().Selections = make(map[string]struct{})

	refreshSelections(ctx)

	return nil
}

func SwitchMode(ctx context.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParameter
	}

	return ctx.PushMode(params[0].(string))
}

func PopMode(ctx context.Context, _ ...interface{}) error {
	return ctx.PopMode()
}

func Refresh(ctx context.Context, params ...interface{}) error {
	currentNode := fs.GetFileManager().Dir.VisibleNodes[ctx.State().FocusIdx]

	focus := currentNode.AbsolutePath

	if len(params) == 1 {
		forcusPath, ok := params[0].(string)

		if ok {
			focus = forcusPath
		}
	}

	ChangeDirectory(ctx, fs.GetFileManager().Dir.Path, false, &focus)

	return nil
}

func refreshSelections(ctx context.Context) {
	gui.GetGui().Views.Selection.SetTitle(len(ctx.State().Selections))
	gui.GetGui().Views.Selection.RenderSelections(ctx.State().Selections)

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)
}

func Quit(_ context.Context, _ ...interface{}) error {
	return gocui.ErrQuit
}
