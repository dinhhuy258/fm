package gui

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/gocui"
)

const (
	horizontalMargin = 1
	verticalMargin   = 1
	logSize          = 2
)

type viewDimension struct {
	x0, y0, x1, y1 int
}

func (gui *Gui) setViewDimensions() error {
	width, height := gui.g.Size()
	width--
	height--

	viewNameMappings := []struct {
		name      string
		dimension viewDimension
	}{
		{
			name: "explorer-header",
			dimension: viewDimension{
				x0: 0,
				y0: 0,
				x1: int(float32(width)*0.7) - horizontalMargin,
				y1: height - logSize - verticalMargin,
			},
		},
		{
			name: "explorer",
			dimension: viewDimension{
				x0: 0,
				y0: verticalMargin,
				x1: int(float32(width)*0.7) - horizontalMargin,
				y1: height - logSize - verticalMargin,
			},
		},
		{
			name: "log",
			dimension: viewDimension{
				x0: 0,
				y0: height - logSize,
				x1: int(float32(width)*0.7) - horizontalMargin,
				y1: height,
			},
		},
		{
			name: "input",
			dimension: viewDimension{
				x0: 0,
				y0: height - logSize,
				x1: int(float32(width)*0.7) - horizontalMargin,
				y1: height,
			},
		},
		{
			name: "selection",
			dimension: viewDimension{
				x0: int(float32(width) * 0.7),
				y0: 0,
				x1: width,
				y1: height/2 - verticalMargin,
			},
		},
		{
			name: "help",
			dimension: viewDimension{
				x0: int(float32(width) * 0.7),
				y0: height / 2,
				x1: width,
				y1: height,
			},
		},
	}

	for _, mapping := range viewNameMappings {
		dimension := mapping.dimension
		_, err := gui.g.SetView(mapping.name, dimension.x0, dimension.y0, dimension.x1, dimension.y1, 0)

		if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}

	return nil
}

func (gui *Gui) layout(_ *gocui.Gui) error {
	if err := gui.setViewDimensions(); err != nil {
		return err
	}

	if err := gui.views.Layout(); err != nil {
		return err
	}

	// Re-focus explorer on resize
	// FIXME: Is there any better way to do this?
	explorerController, _ := gui.controllers.GetController(controller.Explorer).(*controller.ExplorerController)
	explorerController.Focus()

	return nil
}
