package controller

import (
	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

type Type int8

const (
	Explorer Type = iota
	Help
	Sellection
	Progress
	Log
	Input
)

type Event int8

const (
	InputDone Event = iota
)

type Mediator interface {
	notify(Event, string)
}

type IController interface{}

type BaseController struct {
	IController

	mediator Mediator
}

type Controllers struct {
	controllers map[Type]IController
}

func CreateAllControllers(views *view.Views) *Controllers {
	// Selections object to share between explorer and selection controllers
	selections := set.NewSet[string]()
	c := &Controllers{}

	baseController := &BaseController{
		mediator: c,
	}

	c.controllers = make(map[Type]IController)
	c.controllers[Explorer] = newExplorerController(baseController, views.Explorer, selections)
	c.controllers[Sellection] = newSelectionController(baseController, views.Selection, selections)
	c.controllers[Help] = newHelpController(baseController, views.Help)
	c.controllers[Progress] = newProgressController(baseController, views.Progress)
	c.controllers[Log] = newLogController(baseController, views.Log)
	c.controllers[Input] = newInputController(baseController, views.Input)

	return c
}

func (c *Controllers) GetController(controllerType Type) IController {
	return c.controllers[controllerType]
}

func (c *Controllers) notify(event Event, data string) {
	explorerController, _ := c.controllers[Explorer].(*ExplorerController)
	logController, _ := c.controllers[Log].(*LogController)

	if event == InputDone {
		explorerController.view.SetAsCurrentView()
		logController.view.SetViewOnTop()
	}
}
