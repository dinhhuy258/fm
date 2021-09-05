package gui

import (
	"errors"

	"github.com/jroimartin/gocui"
)

type viewDimension struct {
	x0, y0, x1, y1 int
}

func (gui *Gui) createAllViews() error {
	viewNameMappings := []struct {
		viewPtr **gocui.View
		name    string
	}{
		{viewPtr: &gui.Views.Main, name: "main"},
		{viewPtr: &gui.Views.Selection, name: "selection"},
		{viewPtr: &gui.Views.SortAndFilter, name: "sortAndFilter"},
		{viewPtr: &gui.Views.HelpMenu, name: "helpMenu"},
		{viewPtr: &gui.Views.InputAndLogs, name: "inputAndLogs"},
	}

	var err error
	for _, mapping := range viewNameMappings {
		*mapping.viewPtr, err = gui.g.SetView(mapping.name, 0, 0, 10, 10)
		if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}

	gui.Views.Selection.Title = "Main"
	gui.Views.Selection.Title = "Selection"
	gui.Views.SortAndFilter.Title = "Sort & filter"
	gui.Views.HelpMenu.Title = "Help"
	gui.Views.InputAndLogs.Title = "Input"

	return nil
}

func (gui *Gui) setViewDimentions() error {
	width, height := gui.g.Size()
	width--
	height--

	viewNameMappings := []struct {
		name      string
		dimension viewDimension
	}{
		{
			name: "main",
			dimension: viewDimension{
				x0: 0,
				y0: 2,
				x1: int(float32(width) * 0.7),
				y1: height - 2,
			},
		},
		{
			name: "sortAndFilter",
			dimension: viewDimension{
				x0: 0,
				y0: 0,
				x1: int(float32(width) * 0.7),
				y1: 2,
			},
		},
		{
			name: "inputAndLogs",
			dimension: viewDimension{
				x0: 0,
				y0: height - 2,
				x1: int(float32(width) * 0.7),
				y1: height,
			},
		},
		{
			name: "selection",
			dimension: viewDimension{
				x0: int(float32(width) * 0.7),
				y0: 0,
				x1: width - 1,
				y1: height / 2,
			},
		},
		{
			name: "helpMenu",
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
		_, err := gui.g.SetView(mapping.name, dimension.x0, dimension.y0, dimension.x1, dimension.y1)

		if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}

	return nil
}

func (gui *Gui) layout(g *gocui.Gui) error {
	if !gui.ViewsSetup {
		if err := gui.createAllViews(); err != nil {
			return err
		}
	}

	if err := gui.setViewDimentions(); err != nil {
		return err
	}

	return nil
}
