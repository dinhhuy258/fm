package msg

import (
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/key"
)

// SetInputBuffer is a message to set the input buffer
func SetInputBuffer(app IApp, _ key.Key, ctx MessageContext) {
	input := ctx["arg1"].(string)

	inputController, _ := app.GetController(controller.Input).(*controller.InputController)
	inputController.SetInputBuffer(input)
}

// UpdateInputBufferFromKey is a message to update the input buffer from a key
func UpdateInputBufferFromKey(app IApp, key key.Key, _ MessageContext) {
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)
	inputController.UpdateInputBufferFromKey(key)
}
