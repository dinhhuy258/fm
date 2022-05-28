package msg

import (
	"path"
	"strings"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

func NewFileFromInput(app IApp, _ ...string) {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	logController, _ := app.GetController(controller.Log).(*controller.LogController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	name := inputController.GetInputBuffer()
	if name == "" {
		logController.SetLog(view.Warning, "File name is empty")
		logController.UpdateView()

		return
	}

	var err error

	if strings.HasSuffix(name, "/") {
		err = fs.CreateDirectory(name)
	} else {
		err = fs.CreateFile(name, false)
	}

	if err != nil {
		logController.SetLog(view.Error, "Failed to create file %s", name)
		logController.UpdateView()
	} else {
		logController.SetLog(view.Info, "File %s were created successfully", name)
		logController.UpdateView()

		// Reload the current directory in case file were created successfully
		focusPath := path.Join(explorerController.GetPath(), name)
		loadDirectory(app, explorerController.GetPath(), optional.New(focusPath))
	}
}
