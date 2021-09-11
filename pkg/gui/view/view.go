package view

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dinhhuy258/gocui"
)

type Views struct {
	Main          *MainView
	Selection     *SelectionView
	Help          *HelpView
	SortAndFilter *View
	Input         *View
	Log           *View
	Confirm       *View
	Progress      *View
}

func CreateAllViews(g *gocui.Gui) (*Views, error) {
	var (
		mainHeader    *gocui.View
		main          *gocui.View
		selection     *gocui.View
		sortAndFilter *gocui.View
		help          *gocui.View
		input         *gocui.View
		log           *gocui.View
		confirm       *gocui.View
		progress      *gocui.View
	)

	viewNameMappings := []struct {
		viewPtr **gocui.View
		name    string
	}{
		{viewPtr: &mainHeader, name: "main-header"},
		{viewPtr: &main, name: "main"},
		{viewPtr: &selection, name: "selection"},
		{viewPtr: &sortAndFilter, name: "sortAndFilter"},
		{viewPtr: &help, name: "help"},
		{viewPtr: &input, name: "input"},
		{viewPtr: &log, name: "log"},
		{viewPtr: &confirm, name: "confirm"},
		{viewPtr: &progress, name: "progress"},
	}

	var err error
	for _, mapping := range viewNameMappings {
		*mapping.viewPtr, err = g.SetView(mapping.name, 0, 0, 10, 10)
		if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
			return nil, err
		}
	}

	return &Views{
		Main:          newMainView(g, main, mainHeader),
		Selection:     newSelectionView(g, selection),
		SortAndFilter: newView(g, sortAndFilter),
		Help:          newHelpView(g, help),
		Input:         newView(g, input),
		Log:           newView(g, log),
		Confirm:       newView(g, confirm),
		Progress:      newView(g, progress),
	}, nil
}

func (v *Views) Layout() error {
	v.Help.layout()

	return v.Main.layout()
}

type View struct {
	g *gocui.Gui
	v *gocui.View
}

func newView(g *gocui.Gui, v *gocui.View) *View {
	return &View{
		g: g,
		v: v,
	}
}

func (view *View) SetViewContent(displayStrings []string) {
	view.g.Update(func(g *gocui.Gui) error {
		view.v.Clear()
		fmt.Fprint(view.v, strings.Join(displayStrings, "\n"))

		return nil
	})
}

func (view *View) Size() (x, y int) {
	return view.v.Size()
}

func (view *View) SetOrigin(x, y int) error {
	return view.v.SetOrigin(x, y)
}

func (view *View) SetCursor(x, y int) error {
	return view.v.SetCursor(x, y)
}

func (view *View) SetTitle(title string) {
	view.v.Title = title
}

func (view *View) NextCursor() error {
	cx, cy := view.v.Cursor()
	if err := view.v.SetCursor(cx, cy+1); err != nil {
		ox, oy := view.v.Origin()
		if err := view.v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	return nil
}

func (view *View) PreviousCursor() error {
	cx, cy := view.v.Cursor()
	if err := view.v.SetCursor(cx, cy-1); err != nil {
		ox, oy := view.v.Origin()
		if err := view.v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}

	return nil
}
