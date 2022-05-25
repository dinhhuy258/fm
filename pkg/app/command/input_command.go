package command

import "github.com/dinhhuy258/fm/pkg/gui/controller"

func SetInputBuffer(app IApp, params ...string) {
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	input := params[0]
	inputController.SetInputBuffer(input)
}

func UpdateInputBufferFromKey(app IApp, params ...string) {
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	inputController.UpdateInputBufferFromKey(app.GetPressedKey())
}
