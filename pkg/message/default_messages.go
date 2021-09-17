package message

import (
	"errors"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/gocui"
)

var ErrInvalidMessageParams = errors.New("invalid message params")

func ToggleSelection(ctx *ctx.Context, params ...interface{}) error {
	path := (*ctx).GetFileManager().Dir.VisibleNodes[(*ctx).GetState().FocusIdx].AbsolutePath

	if _, hasPath := (*ctx).GetState().Selections[path]; hasPath {
		delete((*ctx).GetState().Selections, path)
	} else {
		(*ctx).GetState().Selections[path] = struct{}{}
	}

	(*ctx).GetGui().Views.Selection.SetTitle(len((*ctx).GetState().Selections))

	if err := (*ctx).GetGui().Views.Selection.RenderSelections((*ctx).GetState().Selections); err != nil {
		return err
	}

	return (*ctx).GetGui().Views.Main.RenderDir(
		(*ctx).GetFileManager().Dir,
		(*ctx).GetState().Selections,
		(*ctx).GetState().FocusIdx,
	)
}

func ToggleHidden(ctx *ctx.Context, params ...interface{}) error {
	config.AppConfig.ShowHidden = !config.AppConfig.ShowHidden

	(*ctx).GetFileManager().Dir.Reload()

	nodeSize := len((*ctx).GetFileManager().Dir.VisibleNodes)
	(*ctx).GetGui().Views.Main.SetTitle(" " + (*ctx).GetFileManager().Dir.Path + " (" + strconv.Itoa(nodeSize) + ") ")
	(*ctx).GetGui().Views.SortAndFilter.SetSortAndFilter()

	return (*ctx).GetGui().Views.Main.RenderDir(
		(*ctx).GetFileManager().Dir,
		(*ctx).GetState().Selections,
		(*ctx).GetState().FocusIdx,
	)
}

func ClearSelection(ctx *ctx.Context, params ...interface{}) error {
	(*ctx).GetState().Selections = make(map[string]struct{})

	(*ctx).GetGui().Views.Selection.SetTitle(len((*ctx).GetState().Selections))

	if err := (*ctx).GetGui().Views.Selection.RenderSelections((*ctx).GetState().Selections); err != nil {
		return err
	}

	return (*ctx).GetGui().Views.Main.RenderDir(
		(*ctx).GetFileManager().Dir,
		(*ctx).GetState().Selections,
		(*ctx).GetState().FocusIdx,
	)
}

func Focus(ctx *ctx.Context, path string) error {
	count := 0

	for _, node := range (*ctx).GetFileManager().Dir.VisibleNodes {
		if node.IsDir && node.AbsolutePath == path {
			break
		}

		count++
	}

	if count == len((*ctx).GetFileManager().Dir.VisibleNodes) {
		return nil
	}

	for i := 0; i < count; i++ {
		if err := (*ctx).GetGui().Views.Main.NextCursor(); err != nil {
			return err
		}

		(*ctx).GetState().FocusIdx++
	}

	return nil
}

func SwitchMode(ctx *ctx.Context, params ...interface{}) error {
	if len(params) != 1 {
		return ErrInvalidMessageParams
	}

	return (*ctx).PushMode(params[0].(string))
}

func PopMode(ctx *ctx.Context, params ...interface{}) error {
	return (*ctx).PopMode()
}

func Refresh(ctx *ctx.Context, params ...interface{}) error {
	(*ctx).GetFileManager().Reload()

	return nil
}

func Quit(ctx *ctx.Context, params ...interface{}) error {
	return gocui.ErrQuit
}
