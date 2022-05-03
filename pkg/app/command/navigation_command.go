package command

import (
	"path/filepath"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

func FocusFirst(app IApp, _ ...interface{}) error {
	gui.GetGui().ResetCursor()
	app.SetFocusIdx(0)

	gui.GetGui().RenderDir(
		fs.GetFileManager().Dir,
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusNext(app IApp, _ ...interface{}) error {
	focusIdx := app.GetFocusIdx()
	if focusIdx == app.GetNumberOfFiles()-1 {
		return nil
	}

	gui.GetGui().NextCursor()
	app.SetFocusIdx(focusIdx + 1)

	gui.GetGui().RenderDir(
		fs.GetFileManager().Dir,
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusPrevious(app IApp, _ ...interface{}) error {
	focusIdx := app.GetFocusIdx()
	if focusIdx == 0 {
		return nil
	}

	gui.GetGui().PreviousCursor()
	app.SetFocusIdx(focusIdx - 1)

	gui.GetGui().RenderDir(
		fs.GetFileManager().Dir,
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusPath(app IApp, params ...interface{}) error {
	path, _ := params[0].(string)
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

	gui.GetGui().ResetCursor()
	app.SetFocusIdx(0)

	for i := 0; i < focusIdx; i++ {
		gui.GetGui().NextCursor()
		app.SetFocusIdx(app.GetFocusIdx() + 1)
	}

	gui.GetGui().RenderDir(
		fs.GetFileManager().Dir,
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func Enter(app IApp, _ ...interface{}) error {
	currentNode := fs.GetFileManager().Dir.VisibleNodes[app.GetFocusIdx()]

	if currentNode.IsDir {
		ChangeDirectory(app, currentNode.AbsolutePath, true, nil)
	}

	return nil
}

func Back(app IApp, _ ...interface{}) error {
	parent := fs.GetFileManager().Dir.Parent()

	ChangeDirectory(app, parent, true, &fs.GetFileManager().Dir.Path)

	return nil
}

func LastVisitedPath(app IApp, _ ...interface{}) error {
	app.VisitLastHistory()
	node := app.PeekHistory()
	ChangeDirectory(app, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func NextVisitedPath(app IApp, _ ...interface{}) error {
	app.VisitNextHistory()
	node := app.PeekHistory()
	ChangeDirectory(app, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func ChangeDirectory(app IApp, path string, saveHistory bool, focusPath *string) {
	if saveHistory && fs.GetFileManager().Dir != nil {
		currentNode := fs.GetFileManager().Dir.VisibleNodes[app.GetFocusIdx()]
		app.PushHistory(currentNode)
	}

	dirLoadedChan := fs.GetFileManager().LoadDirectory(path)

	go func() {
		<-dirLoadedChan

		numberOfFiles := len(fs.GetFileManager().Dir.VisibleNodes)
		app.SetNumberOfFiles(numberOfFiles)
		title := (" " + fs.GetFileManager().Dir.Path + " (" + strconv.Itoa(numberOfFiles) + ") ")
		gui.GetGui().SetMainTitle(title)

		if focusPath == nil {
			_ = FocusFirst(app)
		} else {
			_ = FocusPath(app, *focusPath)
		}

		if saveHistory {
			currentNode := fs.GetFileManager().Dir.VisibleNodes[app.GetFocusIdx()]
			app.PushHistory(currentNode)
		}
	}()
}
