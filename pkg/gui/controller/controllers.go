package controller

import (
	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type ControllerEvent int8

const (
	INPUT_DONE ControllerEvent = iota
)

type ControllerMediator interface {
	notify(ControllerEvent, string)
}

type BaseController struct {
	mediator ControllerMediator
}

type Controllers struct {
	Explorer   *ExplorerController
	Help       *HelpController
	Sellection *SelectionController
	Progress   *ProgressController
	Log        *LogController
	Input      *InputController
}

func CreateAllControllers(views *view.Views) *Controllers {
	// Selections object to share between explorer and selection controllers
	selections := set.NewSet[string]()
	controllers := &Controllers{}

	baseController := &BaseController{
		mediator: controllers,
	}

	controllers.Explorer = newExplorerController(baseController, views.Explorer, selections)
	controllers.Sellection = newSelectionController(baseController, views.Selection, selections)
	controllers.Help = newHelpController(baseController, views.Help)
	controllers.Progress = newProgressController(baseController, views.Progress)
	controllers.Log = newLogController(baseController, views.Log)
	controllers.Input = newInputController(baseController, views.Input)

	return controllers
}

func (c *Controllers) notify(event ControllerEvent, data string) {
	switch event {
	case INPUT_DONE:
		c.Explorer.view.SetAsCurrentView()
	}
}
