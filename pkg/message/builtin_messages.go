package message

import (
	"errors"
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/gocui"
)

var ErrInvalidMessageParams = errors.New("invalid message params")

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

func Focus(ctx *ctx.Context, path string) error {
	count := 0

	for _, node := range (*ctx).GetFileManager().Dir.Nodes {
		if node.IsDir && node.AbsolutePath == path {
			break
		}

		count++
	}

	if count == len((*ctx).GetFileManager().Dir.Nodes) {
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

func ToggleSelection(ctx *ctx.Context, params ...interface{}) error {
	path := (*ctx).GetFileManager().Dir.Nodes[(*ctx).GetState().FocusIdx].AbsolutePath

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

func deletePaths(ctx *ctx.Context, paths []string) error {
	if err := PopMode(ctx); err != nil {
		return err
	}

	if err := (*ctx).GetGui().Views.Main.SetAsCurrentView(); err != nil {
		return err
	}

	if err := (*ctx).GetGui().Views.Progress.StartProgress(1); err != nil {
		return err
	}

	(*ctx).GetFileManager().Delete(paths)

	go func() {
		errCount := 0
	loop:
		for {
			select {
			case <-(*ctx).GetFileManager().DeleteCountChan:
				(*ctx).GetGui().Views.Progress.AddCurrent(1)

				break loop
			case <-(*ctx).GetFileManager().DeleteErrChan:
				errCount++
			}
		}

		var err error
		if errCount != 0 {
			err = (*ctx).GetGui().Views.Log.SetLog(
				fmt.Sprintf("Finished to delete %v. Error count: %d", paths, errCount),
				view.LogLevel(view.INFO),
			)
		} else {
			err = (*ctx).GetGui().Views.Log.SetLog(fmt.Sprintf("Finished to delete file %v", paths), view.LogLevel(view.INFO))
		}

		if err != nil {
			log.Fatalf("failed to set log %v", err)
		}

		if err := Refresh(ctx); err != nil {
			log.Fatalf("failed to refresh %v", err)
		}
	}()

	return nil
}

func DeleteSelections(ctx *ctx.Context, params ...interface{}) error {
	if len((*ctx).GetState().Selections) == 0 {
		err := (*ctx).GetGui().Views.Log.SetLog("Select nothing!!!", view.LogLevel(view.WARNING))
		if err != nil {
			return err
		}

		return PopMode(ctx)
	}

	onYes := func() {
		paths := make([]string, 0, len((*ctx).GetState().Selections))
		for k := range (*ctx).GetState().Selections {
			paths = append(paths, k)
		}

		if err := deletePaths(ctx, paths); err != nil {
			log.Fatalf("failed to delete paths log %v", err)
		}

		// Clear selections
		for k := range (*ctx).GetState().Selections {
			delete((*ctx).GetState().Selections, k)
		}
	}

	onNo := func() {
		if err := PopMode(ctx); err != nil {
			log.Fatalf("failed to pop mode %v", err)
		}

		if err := (*ctx).GetGui().Views.Main.SetAsCurrentView(); err != nil {
			log.Fatalf("failed to set main as the current view %v", err)
		}

		if err := (*ctx).GetGui().Views.Log.SetLog(
			"Canceled deleting the current file/folder",
			view.LogLevel(view.WARNING)); err != nil {
			log.Fatalf("failed to set log %v", err)
		}
	}

	return (*ctx).GetGui().Views.Confirm.SetConfirmation(
		"Do you want to delete selected paths?",
		onYes,
		onNo,
	)
}

func DeleteCurrent(ctx *ctx.Context, params ...interface{}) error {
	currentNode := (*ctx).GetFileManager().Dir.Nodes[(*ctx).GetState().FocusIdx]

	onYes := func() {
		if err := deletePaths(ctx, []string{currentNode.AbsolutePath}); err != nil {
			log.Fatalf("failed to delete paths log %v", err)
		}
	}

	onNo := func() {
		if err := PopMode(ctx); err != nil {
			log.Fatalf("failed to pop mode %v", err)
		}

		if err := (*ctx).GetGui().Views.Main.SetAsCurrentView(); err != nil {
			log.Fatalf("failed to set main as the current view %v", err)
		}

		if err := (*ctx).GetGui().Views.Log.SetLog("Canceled deleting the current file/folder",
			view.LogLevel(view.WARNING)); err != nil {
			log.Fatalf("failed to set log %v", err)
		}
	}

	return (*ctx).GetGui().Views.Confirm.SetConfirmation(
		"Do you want to delete "+currentNode.RelativePath+"?",
		onYes,
		onNo,
	)
}

func Quit(ctx *ctx.Context, params ...interface{}) error {
	return gocui.ErrQuit
}

func changeDirectory(ctx *ctx.Context, path string, saveHistory bool) {
	if saveHistory {
		(*ctx).GetState().History.Push((*ctx).GetFileManager().Dir.Path)
	}

	(*ctx).GetFileManager().LoadDirectory(path)
}
