package controller

import (
	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/fm/pkg/optional"
	"github.com/dinhhuy258/fm/pkg/pipe"
	"github.com/dinhhuy258/gocui"
)

type Type int8

const (
	Explorer Type = iota
	Help
	Sellection
	Log
	Input
	Mark
)

type Event int8

const (
	ShowErrorLog Event = iota
	LogHidden
	CursorEnabled
	CursorDisabled
)

// Mediator is the mediator between controllers
type Mediator interface {
	notify(Event, optional.Optional[string])
}

// IController interface
type IController interface{}

// BaseController is the base controller for all controllers
type BaseController struct {
	IController

	mediator Mediator
}

// Controllers contains all controllers
type Controllers struct {
	g           *gocui.Gui
	controllers map[Type]IController
}

// CreateControllers creates controllers object
func CreateControllers(g *gocui.Gui, views *view.Views, pipe *pipe.Pipe) *Controllers {
	// Selections object to share between explorer and selection controllers
	selections := set.NewSet[string]()
	c := &Controllers{
		g: g,
	}

	baseController := &BaseController{
		mediator: c,
	}

	c.controllers = make(map[Type]IController)
	c.controllers[Explorer] = newExplorerController(baseController, views.Explorer,
		views.ExplorerHeader, selections)
	c.controllers[Sellection] = newSelectionController(baseController, views.Selection, selections)
	c.controllers[Help] = newHelpController(baseController, views.Help)
	c.controllers[Log] = newLogController(baseController, views.Log)
	c.controllers[Input] = newInputController(baseController, views.Input)
	c.controllers[Mark] = newMarkController(baseController, pipe)

	return c
}

// GetController returns controller by type
func (c *Controllers) GetController(controllerType Type) IController {
	return c.controllers[controllerType]
}

// notify is a method that receives events from other controllers
func (c *Controllers) notify(event Event, data optional.Optional[string]) {
	logController, _ := c.controllers[Log].(*LogController)

	switch event {
	case ShowErrorLog:
		data.IfPresent(func(logMsg *string) {
			logController.SetLog(view.Error, *logMsg)
		})
	case LogHidden:
		logController.SetVisible(false)
	case CursorEnabled:
		c.g.Cursor = true
	case CursorDisabled:
		c.g.Cursor = false
	}
}
