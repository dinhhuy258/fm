package message

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/gocui"
)

var ErrInvalidMessageParameter = errors.New("invalid message parameter")

func ToggleSelection(ctx ctx.Context, _ ...interface{}) error {
	path := ctx.FileManager().Dir.VisibleNodes[ctx.State().FocusIdx].AbsolutePath

	if _, hasPath := ctx.State().Selections[path]; hasPath {
		delete(ctx.State().Selections, path)
	} else {
		ctx.State().Selections[path] = struct{}{}
	}

	refreshSelections(ctx)

	return nil
}

func ToggleHidden(ctx ctx.Context, _ ...interface{}) error {
	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	ctx.FileManager().Dir.Reload()

	numberOfFiles := len(ctx.FileManager().Dir.VisibleNodes)
	ctx.State().NumberOfFiles = numberOfFiles
	ctx.Gui().Views.Main.SetTitle(ctx.FileManager().Dir.Path, numberOfFiles)
	ctx.Gui().Views.SortAndFilter.SetSortAndFilter()

	ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)

	return nil
}

func ClearSelection(ctx ctx.Context, _ ...interface{}) error {
	ctx.State().Selections = make(map[string]struct{})

	refreshSelections(ctx)

	return nil
}

func SwitchMode(ctx ctx.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParameter
	}

	return ctx.PushMode(params[0].(string))
}

func PopMode(ctx ctx.Context, _ ...interface{}) error {
	return ctx.PopMode()
}

func Refresh(ctx ctx.Context, params ...interface{}) error {
	currentNode := ctx.FileManager().Dir.VisibleNodes[ctx.State().FocusIdx]

	focus := currentNode.AbsolutePath

	if len(params) == 1 {
		forcusPath, ok := params[0].(string)

		if ok {
			focus = forcusPath
		}
	}

	ChangeDirectory(ctx, ctx.FileManager().Dir.Path, false, &focus)
	ctx.FileManager().LoadDirectory(ctx.FileManager().Dir.Path)

	return nil
}

func refreshSelections(ctx ctx.Context) {
	ctx.Gui().Views.Selection.SetTitle(len(ctx.State().Selections))
	ctx.Gui().Views.Selection.RenderSelections(ctx.State().Selections)

	ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)
}

func Quit(_ ctx.Context, _ ...interface{}) error {
	return gocui.ErrQuit
}
