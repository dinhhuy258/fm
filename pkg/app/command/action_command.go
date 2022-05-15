package command

import (
	"path"
	"strings"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func NewFile(app IApp, _ ...interface{}) error {
	appGui := app.GetGui()
	logController := appGui.GetControllers().Log
	explorerController := appGui.GetControllers().Explorer
	inputController := appGui.GetControllers().Input

	inputController.SetInput(controller.Input, "new file", func(name string) {
		if name == "" {
			logController.SetLog(view.Warning, "File name is empty")

			return
		}

		var err error

		if strings.HasSuffix(name, "/") {
			err = fs.CreateDirectory(name)
		} else {
			err = fs.CreateFile(name)
		}

		if err != nil {
			logController.SetLog(view.Error, "Failed to create file %s", name)
		} else {
			logController.SetLog(view.Info, "File %s were created successfully", name)
			// Reload the current directory in case file were created successfully
			loadDirectory(app, explorerController.GetPath(), path.Join(explorerController.GetPath(), name))
		}
	})

	return nil
}
