package gui

import (
	"errors"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/gocui"
)

const (
	horizontalMargin  = 1
	verticalMargin    = 1
	sortAndFilterSize = 2
	logSize           = 2
)

type viewDimension struct {
	x0, y0, x1, y1 int
}

func (gui *Gui) createAllViews() error {
	viewNameMappings := []struct {
		viewPtr **gocui.View
		name    string
	}{
		{viewPtr: &gui.Views.MainHeader, name: "main-header"},
		{viewPtr: &gui.Views.Main, name: "main"},
		{viewPtr: &gui.Views.Selection, name: "selection"},
		{viewPtr: &gui.Views.SortAndFilter, name: "sortAndFilter"},
		{viewPtr: &gui.Views.HelpMenu, name: "helpMenu"},
		{viewPtr: &gui.Views.Input, name: "input"},
		{viewPtr: &gui.Views.Log, name: "log"},
		{viewPtr: &gui.Views.Confirm, name: "confirm"},
		{viewPtr: &gui.Views.Progress, name: "progress"},
	}

	var err error
	for _, mapping := range viewNameMappings {
		*mapping.viewPtr, err = gui.g.SetView(mapping.name, 0, 0, 10, 10)
		if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}

	gui.initMainPanels()
	gui.initSelectionPanels()

	gui.Views.SortAndFilter.Title = " Sort & filter "

	gui.Views.HelpMenu.Title = " Help "

	gui.Views.Log.Title = " Logs "

	if _, err := gui.g.SetViewOnTop(gui.Views.Log.Name()); err != nil {
		return err
	}

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
			name: "main-header",
			dimension: viewDimension{
				x0: 0,
				y0: sortAndFilterSize + verticalMargin,
				x1: int(float32(width)*0.7) - horizontalMargin,
				y1: height - logSize - verticalMargin,
			},
		},
		{
			name: "main",
			dimension: viewDimension{
				x0: 0,
				y0: sortAndFilterSize + verticalMargin + 1, // plus 1 for header
				x1: int(float32(width)*0.7) - horizontalMargin,
				y1: height - logSize - verticalMargin,
			},
		},
		{
			name: "sortAndFilter",
			dimension: viewDimension{
				x0: 0,
				y0: 0,
				x1: int(float32(width)*0.7) - horizontalMargin,
				y1: sortAndFilterSize,
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
			name: "confirm",
			dimension: viewDimension{
				x0: 0,
				y0: height - logSize,
				x1: int(float32(width)*0.7) - horizontalMargin,
				y1: height,
			},
		},
		{
			name: "progress",
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

		gui.GuiLoadedChan <- struct{}{}
		close(gui.GuiLoadedChan)
		gui.ViewsSetup = true
	}

	if err := gui.setViewDimentions(); err != nil {
		return err
	}

	x, _ := gui.Views.Main.Size()
	gui.MainRow.SetWidth(x)

	rowString, err := gui.MainRow.HeaderRow.Sprint(
		[]string{config.AppConfig.IndexHeader, config.AppConfig.PathHeader, config.AppConfig.SizeHeader},
	)
	if err != nil {
		return err
	}

	gui.SetViewContent(gui.Views.MainHeader, []string{rowString})

	return nil
}
