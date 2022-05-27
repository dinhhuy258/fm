package controller

import (
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/key"
	"github.com/dinhhuy258/fm/pkg/optional"
)

type InputController struct {
	*BaseController

	view *view.InputView
}

func newInputController(baseController *BaseController, view *view.InputView) *InputController {
	inputController := &InputController{
		BaseController: baseController,
		view:           view,
	}

	return inputController
}

func (ic *InputController) SetInputBuffer(input string) {
	ic.view.SetInputBuffer(input)
	ic.mediator.notify(LogHidden, optional.NewEmpty[string]())

}

func (ic *InputController) GetInputBuffer() string {
	return ic.view.GetInputBuffer()
}

func (ic *InputController) UpdateInputBufferFromKey(key key.Key) {
	ic.view.UpdateInputBufferFromKey(key)
}
