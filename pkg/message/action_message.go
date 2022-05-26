package message

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
		err = fs.CreateFile(name)
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

func Rename(app IApp, _ ...string) {
	// explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	// logController, _ := app.GetController(controller.Log).(*controller.LogController)
	// inputController, _ := app.GetController(controller.Input).(*controller.InputController)
	//
	// currentEntry := explorerController.GetCurrentEntry()
	//
	// inputController.SetInput(
	// 	controller.InputText,
	// 	"rename",
	// 	optional.New(currentEntry.GetName()),
	// 	func(newName string) {
	// 		if newName == "" {
	// 			logController.SetLog(view.Warning, "File name is empty")
	// 			logController.UpdateView()
	//
	// 			return
	// 		}
	//
	// 		err := fs.Rename(currentEntry.GetPath(), path.Join(explorerController.GetPath(), newName))
	//
	// 		if err != nil {
	// 			logController.SetLog(view.Error, "Failed to rename file %s", newName)
	// 			logController.UpdateView()
	// 		} else {
	// 			logController.SetLog(view.Info, "File %s were renamed successfully", newName)
	// 			logController.UpdateView()
	//
	// 			// Reload the current directory
	// 			focusPath := path.Join(explorerController.GetPath(), newName)
	// 			loadDirectory(app, explorerController.GetPath(), optional.New(focusPath))
	// 		}
	// 	},
	// )
	// inputController.UpdateView()
}
