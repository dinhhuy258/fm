package command

import (
	"path/filepath"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

func FocusFirst(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()

	gui.GetGui().ResetCursor()
	app.SetFocusIdx(0)

	gui.GetGui().RenderEntries(
		fileExplorer.GetEntries(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusNext(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()

	focusIdx := app.GetFocusIdx()
	if focusIdx == fileExplorer.GetEntriesSize()-1 {
		return nil
	}

	gui.GetGui().NextCursor()
	app.SetFocusIdx(focusIdx + 1)

	gui.GetGui().RenderEntries(
		fileExplorer.GetEntries(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusPrevious(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()

	focusIdx := app.GetFocusIdx()
	if focusIdx == 0 {
		return nil
	}

	gui.GetGui().PreviousCursor()
	app.SetFocusIdx(focusIdx - 1)

	gui.GetGui().RenderEntries(
		fileExplorer.GetEntries(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusPath(app IApp, params ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()

	path, _ := params[0].(string)

	if fileExplorer.GetPath() != filepath.Dir(path) {
		fileExplorer.LoadEntries(filepath.Dir(path), func() {
			focusPath(app, path)
		})
	} else {
		focusPath(app, path)
	}

	return nil
}

func Enter(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	entry := fileExplorer.GetEntry(app.GetFocusIdx())

	if entry.IsDirectory() {
		LoadDirectory(app, entry.GetPath(), true, "")
	}

	return nil
}

func Back(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	dir := fileExplorer.Dir()
	if dir == "." {
		// If folder has no parent directory then do nothing
		return nil
	}

	LoadDirectory(app, dir, true, fileExplorer.GetPath())

	return nil
}

func LastVisitedPath(app IApp, _ ...interface{}) error {
	app.VisitLastHistory()
	node := app.PeekHistory()
	// TODO: Move filepath.Dir to fs package
	ChangeDirectory(app, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

func NextVisitedPath(app IApp, _ ...interface{}) error {
	app.VisitNextHistory()
	node := app.PeekHistory()
	// TODO: Move filepath.Dir to fs package
	ChangeDirectory(app, filepath.Dir(node.AbsolutePath), false, &node.AbsolutePath)

	return nil
}

// TODO: Do we need pointer in focusPath variable
func ChangeDirectory(app IApp, path string, saveHistory bool, focusPath *string) {
	fileManager := fs.GetFileManager()
	// TODO Better way to handle isLoaded
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

func LoadDirectory(app IApp, path string, saveHistory bool, focusPath string) {
	fileExplorer := fs.GetFileExplorer()

	fileExplorer.LoadEntries(path, func() {
		entriesSize := fileExplorer.GetEntriesSize()
		path := fileExplorer.GetPath()

		title := (" " + path + " (" + strconv.Itoa(entriesSize) + ") ")
		gui.GetGui().SetMainTitle(title)

		if focusPath == "" {
			_ = FocusFirst(app)
		} else {
			_ = FocusPath(app, focusPath)
		}

		// if saveHistory {
		// currentNode := fileManager.GetNodeAtIdx(app.GetFocusIdx())
		// app.PushHistory(currentNode)
		// }
	})
}

func focusPath(app IApp, path string) {
	fileExplorer := fs.GetFileExplorer()
	focusIdx := 0

	// Iterate through the list of entries and find the idx for the current path
	for idx, entry := range fileExplorer.GetEntries() {
		if entry.GetPath() == path {
			focusIdx = idx
			break
		}
	}

	gui.GetGui().ResetCursor()
	app.SetFocusIdx(0)

	for i := 0; i < focusIdx; i++ {
		gui.GetGui().NextCursor()
		app.SetFocusIdx(app.GetFocusIdx() + 1)
	}

	gui.GetGui().RenderEntries(
		fileExplorer.GetEntries(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)
}
