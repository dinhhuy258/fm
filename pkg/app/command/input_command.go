package command

import "github.com/dinhhuy258/fm/pkg/gui/controller"

func SetInputBuffer(app IApp, params ...interface{}) {
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	input, _ := params[0].(string)
	inputController.SetInputBuffer(input)
}

func UpdateInputBufferFromKey(app IApp, params ...interface{}) {
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	key, _ := params[0].(string)
	inputController.UpdateInputBufferFromKey(key)
}

