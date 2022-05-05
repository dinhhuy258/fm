package command

import (
	"fmt"
	"path"
	"strings"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func NewFile(app IApp, _ ...interface{}) error {
	fileExplorer := fs.GetFileExplorer()
	appGui := gui.GetGui()

	appGui.SetInput("new file", func(name string) {
		if name == "" {
			appGui.SetLog("File name is empty", view.LogLevel(view.WARNING))

			return
		}

		var err error

		if strings.HasSuffix(name, "/") {
			err = fs.CreateDirectory(name)
		} else {
			err = fs.CreateFile(name)
		}

		if err != nil {
			appGui.SetLog(fmt.Sprintf("Failed to create file %s", name), view.LogLevel(view.ERROR))
		} else {
			appGui.SetLog(fmt.Sprintf("File %s were created successfully", name),
				view.LogLevel(view.INFO))
			// Reload the current directory in case file were created successfully
			LoadDirectory(app, fileExplorer.GetPath(), false, path.Join(fileExplorer.GetPath(), name))
		}
	})

	return nil
}
