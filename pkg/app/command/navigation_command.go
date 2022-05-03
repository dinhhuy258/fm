package command

import (
	"path/filepath"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

func FocusFirst(app IApp, _ ...interface{}) error {
	_ = gui.GetGui().Views.Main.SetOrigin(0, 0)
	_ = gui.GetGui().Views.Main.SetCursor(0, 0)
	app.State().FocusIdx = 0

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		app.State().Selections,
		app.State().FocusIdx,
	)

	return nil
}

func FocusNext(app IApp, _ ...interface{}) error {
	if app.State().FocusIdx == app.State().NumberOfFiles-1 {
		return nil
	}

	if err := gui.GetGui().Views.Main.NextCursor(); err != nil {
		return err
	}

	app.State().FocusIdx++

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		app.State().Selections,
		app.State().FocusIdx,
	)

	return nil
}

func FocusPrevious(app IApp, _ ...interface{}) error {
	if app.State().FocusIdx == 0 {
		return nil
	}

	if err := gui.GetGui().Views.Main.PreviousCursor(); err != nil {
		return err
	}

	app.State().FocusIdx--

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		app.State().Selections,
		app.State().FocusIdx,
	)

	return nil
}

func FocusPath(app IApp, params ...interface{}) error {
	path := params[0].(string)
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

	app.State().FocusIdx = 0

	for i := 0; i < focusIdx; i++ {
		if err := gui.GetGui().Views.Main.NextCursor(); err != nil {
			return err
		}

		app.State().FocusIdx++
	}

	gui.GetGui().Views.Main.RenderDir(
		fs.GetFileManager().Dir,
		app.State().Selections,
		app.State().FocusIdx,
	)

	return nil
}

func Enter(app IApp, _ ...interface{}) error {
	currentNode := fs.GetFileManager().Dir.VisibleNodes[app.State().FocusIdx]

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
	app.State().History.VisitLast()
	node := app.State().History.Peek()
	ChangeDirectory(app, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func NextVisitedPath(app IApp, _ ...interface{}) error {
	app.State().History.VisitNext()
	node := app.State().History.Peek()
	ChangeDirectory(app, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func ChangeDirectory(app IApp, path string, saveHistory bool, focusPath *string) {
	if saveHistory && fs.GetFileManager().Dir != nil {
		currentNode := fs.GetFileManager().Dir.VisibleNodes[app.State().FocusIdx]
		app.State().History.Push(currentNode)
	}

	dirLoadedChan := fs.GetFileManager().LoadDirectory(path)

	go func() {
		<-dirLoadedChan

		numberOfFiles := len(fs.GetFileManager().Dir.VisibleNodes)
		app.State().NumberOfFiles = numberOfFiles
		gui.GetGui().Views.Main.SetTitle(fs.GetFileManager().Dir.Path, numberOfFiles)

		if focusPath == nil {
			_ = FocusFirst(app)
		} else {
			_ = FocusPath(app, *focusPath)
		}

		if saveHistory {
			currentNode := fs.GetFileManager().Dir.VisibleNodes[app.State().FocusIdx]
			app.State().History.Push(currentNode)
		}
	}()
}
