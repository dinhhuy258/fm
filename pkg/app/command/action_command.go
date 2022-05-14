package command

import (
	"path"
	"strings"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func NewFile(app IApp, _ ...interface{}) error {
	appGui := gui.GetGui()
	logController := appGui.GetControllers().Log
	fileExplorer := fs.GetFileExplorer()

	appGui.SetInput("new file", func(name string) {
		if name == "" {
			logController.SetLog(view.LogLevel(view.WARNING), "File name is empty")

			return
		}

		var err error

		if strings.HasSuffix(name, "/") {
			err = fs.CreateDirectory(name)
		} else {
			err = fs.CreateFile(name)
		}

		if err != nil {
			logController.SetLog(view.LogLevel(view.ERROR), "Failed to create file %s", name)
		} else {
			logController.SetLog(view.LogLevel(view.INFO), "File %s were created successfully", name)
			// Reload the current directory in case file were created successfully
			LoadDirectory(app, fileExplorer.GetPath(), path.Join(fileExplorer.GetPath(), name))
		}
	})

	return nil
}
