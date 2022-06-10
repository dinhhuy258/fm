package controller

import (
	"github.com/dinhhuy258/fm/pkg/gui/key"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
)

// InputController is a controller for input view
type InputController struct {
	*BaseController

	view *view.InputView
}

// newInputController creates a new input controller
func newInputController(baseController *BaseController, view *view.InputView) *InputController {
	inputController := &InputController{
		BaseController: baseController,
		view:           view,
	}

	return inputController
}

// SetInputBuffer sets `input` to the input buffer, this action will show the input view
func (ic *InputController) SetInputBuffer(input string) {
	ic.view.SetInputBuffer(input)
	ic.mediator.notify(InputVisible, optional.NewEmpty[string]())
}

// GetInputBuffer gets content from the input buffer
func (ic *InputController) GetInputBuffer() string {
	return ic.view.GetInputBuffer()
}

// UpdateInputBufferFromKey updates the input buffer from key
func (ic *InputController) UpdateInputBufferFromKey(key key.Key) {
	ic.view.UpdateInputBufferFromKey(key)
}
