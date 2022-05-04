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
		fs.GetFileManager().GetVisibleNodes(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusNext(app IApp, _ ...interface{}) error {
	fileManager := fs.GetFileManager()
	focusIdx := app.GetFocusIdx()
	if focusIdx == fileManager.GetVisibleNodesSize()-1 {
		return nil
	}

	gui.GetGui().NextCursor()
	app.SetFocusIdx(focusIdx + 1)

	gui.GetGui().RenderDir(
		fileManager.GetVisibleNodes(),
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
		fs.GetFileManager().GetVisibleNodes(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusPath(app IApp, params ...interface{}) error {
	fileManager := fs.GetFileManager()
	path, _ := params[0].(string)
	if fileManager.GetCurrentPath() != filepath.Dir(path) {
		dirLoadedChan := fileManager.LoadDirectory(filepath.Dir(path))
		<-dirLoadedChan
	}

	focusIdx := 0

	for _, node := range fileManager.GetVisibleNodes() {
		if node.AbsolutePath == path {
			break
		}

		focusIdx++
	}

	if focusIdx == fileManager.GetVisibleNodesSize() {
		focusIdx = fileManager.GetVisibleNodesSize() - 1
	}

	gui.GetGui().ResetCursor()
	app.SetFocusIdx(0)

	for i := 0; i < focusIdx; i++ {
		gui.GetGui().NextCursor()
		app.SetFocusIdx(app.GetFocusIdx() + 1)
	}

	gui.GetGui().RenderDir(
		fs.GetFileManager().GetVisibleNodes(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func Enter(app IApp, _ ...interface{}) error {
	currentNode := fs.GetFileManager().GetNodeAtIdx(app.GetFocusIdx())

	if currentNode.IsDir {
		ChangeDirectory(app, currentNode.AbsolutePath, true, nil)
	}

	return nil
}

func Back(app IApp, _ ...interface{}) error {
	fileManager := fs.GetFileManager()
	parent := fileManager.GetParentPath()

	currentPath := fileManager.GetCurrentPath()
	ChangeDirectory(app, parent, true, &currentPath)

	return nil
}

func LastVisitedPath(app IApp, _ ...interface{}) error {
	app.VisitLastHistory()
	node := app.PeekHistory()
	//TODO: Move filepath.Dir to fs package
	ChangeDirectory(app, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func NextVisitedPath(app IApp, _ ...interface{}) error {
	app.VisitNextHistory()
	node := app.PeekHistory()
	//TODO: Move filepath.Dir to fs package
	ChangeDirectory(app, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

//TODO: Do we need pointer in focusPath variable
func ChangeDirectory(app IApp, path string, saveHistory bool, focusPath *string) {
	fileManager := fs.GetFileManager()
	//TODO Better way to handle isLoaded
	if saveHistory && fileManager.IsLoaded() {
		currentNode := fileManager.GetNodeAtIdx(app.GetFocusIdx())
		app.PushHistory(currentNode)
	}

	dirLoadedChan := fs.GetFileManager().LoadDirectory(path)

	go func() {
		<-dirLoadedChan

		numberOfFiles := fileManager.GetVisibleNodesSize()
		title := (" " + fileManager.GetCurrentPath() + " (" + strconv.Itoa(numberOfFiles) + ") ")
		gui.GetGui().SetMainTitle(title)

		if focusPath == nil {
			_ = FocusFirst(app)
		} else {
			_ = FocusPath(app, *focusPath)
		}

		if saveHistory {
			currentNode := fileManager.GetNodeAtIdx(app.GetFocusIdx())
			app.PushHistory(currentNode)
		}
	}()
}
