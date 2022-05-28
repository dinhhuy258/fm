package controller

import (
	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
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
	ShowErrorLog Event = iota
	LogHidden
)

type Mediator interface {
	notify(Event, optional.Optional[string])
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
	c.controllers[Explorer] = newExplorerController(baseController, views.Explorer,
		views.ExplorerHeader, selections)
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

func (c *Controllers) notify(event Event, data optional.Optional[string]) {
	logController, _ := c.controllers[Log].(*LogController)

	switch event {
	case ShowErrorLog:
		data.IfPresent(func(logMsg *string) {
			logController.SetLog(view.Error, *logMsg)
		})
	case LogHidden:
		logController.SetVisible(false)
	}
}
