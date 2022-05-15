package controller

import (
	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type Event int8

const (
	InputDone Event = iota
)

type Mediator interface {
	notify(Event, string)
}

type BaseController struct {
	mediator Mediator
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

func (c *Controllers) notify(event Event, data string) {
	if event == InputDone {
		c.Explorer.view.SetAsCurrentView()
	}
}
