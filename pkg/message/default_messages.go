package message

import (
	"errors"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/gocui"
)

var ErrInvalidMessageParams = errors.New("invalid message params")

func ToggleSelection(ctx ctx.Context, params ...interface{}) error {
	path := ctx.FileManager().Dir.VisibleNodes[ctx.State().FocusIdx].AbsolutePath

	if _, hasPath := ctx.State().Selections[path]; hasPath {
		delete(ctx.State().Selections, path)
	} else {
		ctx.State().Selections[path] = struct{}{}
	}

	ctx.Gui().Views.Selection.SetTitle(len(ctx.State().Selections))

	if err := ctx.Gui().Views.Selection.RenderSelections(ctx.State().Selections); err != nil {
		return err
	}

	return ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)
}

func ToggleHidden(ctx ctx.Context, params ...interface{}) error {
	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	ctx.FileManager().Dir.Reload()

	nodeSize := len(ctx.FileManager().Dir.VisibleNodes)
	ctx.Gui().Views.Main.SetTitle(" " + ctx.FileManager().Dir.Path + " (" + strconv.Itoa(nodeSize) + ") ")
	ctx.Gui().Views.SortAndFilter.SetSortAndFilter()

	return ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)
}

func ClearSelection(ctx ctx.Context, params ...interface{}) error {
	ctx.State().Selections = make(map[string]struct{})

	ctx.Gui().Views.Selection.SetTitle(len(ctx.State().Selections))

	if err := ctx.Gui().Views.Selection.RenderSelections(ctx.State().Selections); err != nil {
		return err
	}

	return ctx.Gui().Views.Main.RenderDir(
		ctx.FileManager().Dir,
		ctx.State().Selections,
		ctx.State().FocusIdx,
	)
}

func Focus(ctx ctx.Context, path string) error {
	count := 0

	for _, node := range ctx.FileManager().Dir.VisibleNodes {
		if node.IsDir && node.AbsolutePath == path {
			break
		}

		count++
	}

	if count == len(ctx.FileManager().Dir.VisibleNodes) {
		return nil
	}

	for i := 0; i < count; i++ {
		if err := ctx.Gui().Views.Main.NextCursor(); err != nil {
			return err
		}

		ctx.State().FocusIdx++
	}

	return nil
}

func SwitchMode(ctx ctx.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParams
	}

	return ctx.PushMode(params[0].(string))
}

func PopMode(ctx ctx.Context, params ...interface{}) error {
	return ctx.PopMode()
}

func Refresh(ctx ctx.Context, params ...interface{}) error {
	ctx.FileManager().Reload()

	return nil
}

func Quit(ctx ctx.Context, params ...interface{}) error {
	return gocui.ErrQuit
}
