package controller

import set "github.com/deckarep/golang-set/v2"

type ControllerEvent int8

const (
	INPUT_DONE ControllerEvent = iota
)

type ControllerMediator interface {
	notify(ControllerEvent, interface{})
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

func CreateAllControllers() *Controllers {
	// Selections object to share between explorer and selection controllers
	selections := set.NewSet[string]()
	controllers := &Controllers{}

	baseController := &BaseController{
		mediator: controllers,
	}

	controllers.Explorer = newExplorerController(baseController, selections)
	controllers.Sellection = newSelectionController(baseController, selections)
	controllers.Help = newHelpController(baseController)
	controllers.Progress = newProgressController(baseController)
	controllers.Log = newLogController(baseController)
	controllers.Input = newInputController(baseController)

	return controllers
}

func (c *Controllers) notify(event ControllerEvent, data interface{}) {
	switch event {
	case INPUT_DONE:
	   c.Explorer.view.SetAsCurrentView()
	}
}
