package msg

import "github.com/dinhhuy258/fm/pkg/gui/controller"

// SetInputBuffer is a message to set the input buffer
func SetInputBuffer(app IApp, params ...string) {
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	input := params[0]
	inputController.SetInputBuffer(input)
}

// UpdateInputBufferFromKey is a message to update the input buffer from a key
func UpdateInputBufferFromKey(app IApp, params ...string) {
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)

	inputController.UpdateInputBufferFromKey(app.GetPressedKey())
}
