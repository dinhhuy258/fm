package command

import (
	"path/filepath"
	"strconv"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
)

func FocusFirst(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()

	appGui.ResetCursor()
	app.SetFocus(0)

	app.RenderEntries()

	return nil
}

func FocusNext(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	focus := app.GetFocus()
	if focus == fileExplorer.GetEntriesSize()-1 {
		return nil
	}

	appGui.NextCursor()
	app.SetFocus(focus + 1)

	app.RenderEntries()

	return nil
}

func FocusPrevious(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()

	focus := app.GetFocus()
	if focus == 0 {
		return nil
	}

	appGui.PreviousCursor()
	app.SetFocus(focus - 1)

	app.RenderEntries()

	return nil
}

func FocusPath(app IApp, params ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()

	if path, _ := params[0].(string); fileExplorer.GetPath() != filepath.Dir(path) {
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

	if fileExplorer.GetEntriesSize() <= 0 {
		// In case of folder is empty
		return nil
	}

	entry := fileExplorer.GetEntry(app.GetFocus())
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

func ChangeDirectory(app IApp, params ...interface{}) error {
	directory, _ := params[0].(string)

	LoadDirectory(app, directory, true, "")

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

		if saveHistory && fileExplorer.GetEntriesSize() > 0 {
			entry := fileExplorer.GetEntry(app.GetFocus())
			app.PushHistory(entry)
		}
	})
}

func focusPath(app IApp, path string) {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	focus := 0
	// Iterate through the list of entries and find the idx for the current path
	for idx, entry := range fileExplorer.GetEntries() {
		if entry.GetPath() == path {
			focus = idx

			break
		}
	}

	appGui.ResetCursor()

	for i := 0; i < focus; i++ {
		appGui.NextCursor()
	}

	app.SetFocus(focus)

	app.RenderEntries()
}
