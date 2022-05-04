package command

import (
	"path/filepath"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

func FocusFirst(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	appGui.ResetCursor()
	app.SetFocusIdx(0)

	appGui.RenderEntries(
		fileExplorer.GetEntries(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusNext(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	focusIdx := app.GetFocusIdx()
	if focusIdx == fileExplorer.GetEntriesSize()-1 {
		return nil
	}

	appGui.NextCursor()
	app.SetFocusIdx(focusIdx + 1)

	appGui.RenderEntries(
		fileExplorer.GetEntries(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)

	return nil
}

func FocusPrevious(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	focusIdx := app.GetFocusIdx()
	if focusIdx == 0 {
		return nil
	}

	appGui.PreviousCursor()
	app.SetFocusIdx(focusIdx - 1)

	appGui.RenderEntries(
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
	entry := app.PeekHistory()

	LoadDirectory(app, fs.Dir(entry.GetPath()), false, entry.GetPath())

	return nil
}

func NextVisitedPath(app IApp, _ ...interface{}) error {
	app.VisitNextHistory()
	entry := app.PeekHistory()

	LoadDirectory(app, fs.Dir(entry.GetPath()), false, entry.GetPath())

	return nil
}

func LoadDirectory(app IApp, path string, saveHistory bool, focusPath string) {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	fileExplorer.LoadEntries(path, func() {
		entriesSize := fileExplorer.GetEntriesSize()
		path := fileExplorer.GetPath()

		title := (" " + path + " (" + strconv.Itoa(entriesSize) + ") ")
		appGui.SetMainTitle(title)

		if focusPath == "" {
			_ = FocusFirst(app)
		} else {
			_ = FocusPath(app, focusPath)
		}

		if saveHistory {
			entry := fileExplorer.GetEntry(app.GetFocusIdx())
			app.PushHistory(entry)
		}
	})
}

func focusPath(app IApp, path string) {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	focusIdx := 0
	// Iterate through the list of entries and find the idx for the current path
	for idx, entry := range fileExplorer.GetEntries() {
		if entry.GetPath() == path {
			focusIdx = idx
			break
		}
	}

	appGui.ResetCursor()
	app.SetFocusIdx(0)

	for i := 0; i < focusIdx; i++ {
		appGui.NextCursor()
		app.SetFocusIdx(app.GetFocusIdx() + 1)
	}

	appGui.RenderEntries(
		fileExplorer.GetEntries(),
		app.GetSelections(),
		app.GetFocusIdx(),
	)
}
