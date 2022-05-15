package view

import (
	"fmt"
	"log"
	"strings"

	"github.com/dinhhuy258/gocui"
)

type Views struct {
	Explorer  *ExplorerView
	Selection *SelectionView
	Help      *HelpView
	Input     *InputView
	Log       *LogView
	Progress  *ProgressView
}

func CreateAllViews(g *gocui.Gui) *Views {
	var (
		explorerHeader *gocui.View
		explorer       *gocui.View
		selection      *gocui.View
		help           *gocui.View
		input          *gocui.View
		log            *gocui.View
		progress       *gocui.View
	)

	viewNameMappings := []struct {
		viewPtr **gocui.View
		name    string
	}{
		{viewPtr: &explorerHeader, name: "explorer-header"},
		{viewPtr: &explorer, name: "explorer"},
		{viewPtr: &selection, name: "selection"},
		{viewPtr: &help, name: "help"},
		{viewPtr: &input, name: "input"},
		{viewPtr: &log, name: "log"},
		{viewPtr: &progress, name: "progress"},
	}

	for _, mapping := range viewNameMappings {
		// No need to handle error here, since we are creating views
		*mapping.viewPtr, _ = g.SetView(mapping.name, 0, 0, 10, 10)
	}

	return &Views{
		Explorer:  newExplorerView(g, explorer, explorerHeader),
		Selection: newSelectionView(g, selection),
		Help:      newHelpView(g, help),
		Input:     newInputView(g, input),
		Log:       newLogView(g, log),
		Progress:  newProgressView(g, progress),
	}
}

func (v *Views) Layout() error {
	v.Help.layout()

	if err := v.Explorer.layout(); err != nil {
		return err
	}

	return nil
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
		_, err := fmt.Fprint(view.v, strings.Join(displayStrings, "\n"))

		return err
	})
}

func (view *View) Size() (int, int) {
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

func (view *View) SetViewOnTop() {
	if _, err := view.g.SetViewOnTop(view.v.Name()); err != nil {
		log.Fatalf("failed to set view %s on top. Error: %v", view.v.Name(), err)
	}
}
