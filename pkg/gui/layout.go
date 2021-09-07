package gui

import (
	"github.com/jesseduffield/gocui"
)

const (
	horizontalMargin  = 1
	verticalMargin    = 1
	sortAndFilterSize = 2
	inputAndLogsSize  = 2
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
		*mapping.viewPtr, err = gui.g.SetView(mapping.name, 0, 0, 10, 10, 0)
		if err != nil && err.Error() != gocui.ErrUnknownView.Error() {
			return err
		}
	}

	gui.Views.Main.Title = ""
	gui.Views.Main.Highlight = true
	gui.Views.Main.SelFgColor = gocui.ColorBlue

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
				y0: sortAndFilterSize + verticalMargin,
				x1: int(float32(width)*0.7) - horizontalMargin,
				y1: height - inputAndLogsSize - verticalMargin,
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
			name: "inputAndLogs",
			dimension: viewDimension{
				x0: 0,
				y0: height - inputAndLogsSize,
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
		_, err := gui.g.SetView(mapping.name, dimension.x0, dimension.y0, dimension.x1, dimension.y1, 0)

		if err != nil && err.Error() != gocui.ErrUnknownView.Error() {
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

	return gui.setViewDimentions()
}
