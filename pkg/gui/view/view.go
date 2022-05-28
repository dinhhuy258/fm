package view

import (
	"strings"

	"github.com/dinhhuy258/gocui"
)

type Views struct {
	Explorer       *ExplorerView
	ExplorerHeader *ExplorerHeaderView
	Selection      *SelectionView
	Help           *HelpView
	Input          *InputView
	Log            *LogView
	Progress       *ProgressView
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
		{viewPtr: &progress, name: "progress"},
		{viewPtr: &input, name: "input"},
		{viewPtr: &log, name: "log"},
	}

	for _, mapping := range viewNameMappings {
		// No need to handle error here, since we are creating views
		*mapping.viewPtr, _ = g.SetView(mapping.name, 0, 0, 10, 10, 0)
	}

	return &Views{
		Explorer:       newExplorerView(g, explorer),
		ExplorerHeader: newExplorerHeaderView(g, explorerHeader),
		Selection:      newSelectionView(g, selection),
		Help:           newHelpView(g, help),
		Input:          newInputView(g, input),
		Log:            newLogView(g, log),
		Progress:       newProgressView(g, progress),
	}
}

func (v *Views) Layout() error {
	if err := v.Explorer.layout(); err != nil {
		return err
	}

	if err := v.ExplorerHeader.layout(); err != nil {
		return err
	}

	if err := v.Selection.layout(); err != nil {
		return err
	}

	if err := v.Help.layout(); err != nil {
		return err
	}

	if err := v.Input.layout(); err != nil {
		return err
	}

	if err := v.Log.layout(); err != nil {
		return err
	}

	if err := v.Progress.layout(); err != nil {
		return err
	}

	return nil
}

type View struct {
	*gocui.View

	g *gocui.Gui
}

func newView(g *gocui.Gui, v *gocui.View) *View {
	return &View{
		View: v,
		g: g,
	}
}

func (view *View) SetViewContent(displayStrings []string) {
	view.SetContent(strings.Join(displayStrings, "\n"))
}

func (view *View) layout() error {
	return nil
}
