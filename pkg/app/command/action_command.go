package command

import (
	"path"
	"strings"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func NewFile(app IApp, _ ...interface{}) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	inputController.SetInput(controller.InputText, "new file", func(name string) {
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
			focusPath := path.Join(explorerController.GetPath(), name)
			loadDirectory(app, explorerController.GetPath(), optional.NewOptional(focusPath))
		}
	})
}
