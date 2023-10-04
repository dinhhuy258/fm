package view

import (
	"github.com/dinhhuy258/gocui"
)

type Views struct {
	Explorer       *ExplorerView
	ExplorerHeader *ExplorerHeaderView
	Input          *InputView
	Log            *LogView
}

func CreateViews(g *gocui.Gui) *Views {
	var (
		explorerHeader *gocui.View
		explorer       *gocui.View
		input          *gocui.View
		log            *gocui.View
	)

	viewNameMappings := []struct {
		viewPtr **gocui.View
		name    string
	}{
		{viewPtr: &explorerHeader, name: "explorer-header"},
		{viewPtr: &explorer, name: "explorer"},
		{viewPtr: &input, name: "input"},
		{viewPtr: &log, name: "log"},
	}

	for _, mapping := range viewNameMappings {
		// No need to handle error here, since we are creating views
		*mapping.viewPtr, _ = g.SetView(mapping.name, 0, 0, 10, 10, 0)
	}

	return &Views{
		Explorer:       newExplorerView(explorer),
		ExplorerHeader: newExplorerHeaderView(explorerHeader),
		Input:          newInputView(input),
		Log:            newLogView(log),
	}
}

func (v *Views) Layout() error {
	if err := v.Explorer.layout(); err != nil {
		return err
	}

	if err := v.ExplorerHeader.layout(); err != nil {
		return err
	}

	if err := v.Input.layout(); err != nil {
		return err
	}

	return v.Log.layout()
}

type View struct {
	*gocui.View
}

func newView(v *gocui.View) *View {
	return &View{
		View: v,
	}
}

func (view *View) layout() error {
	return nil
}
