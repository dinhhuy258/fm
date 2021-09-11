package gui

import "fmt"

func (gui *Gui) initSelectionPanels() {
	gui.SetSelectionTitle(0)
}

func (gui *Gui) SetSelectionTitle(selectionsNum int) {
	gui.Views.Selection.Title = fmt.Sprintf(" Selection (%d) ", selectionsNum)
}

func (gui *Gui) RenderSelections(selections map[string]struct{}) error {
	s := make([]string, 0, len(selections))
	for k := range selections {
		s = append(s, k)
	}

	gui.SetViewContent(gui.Views.Selection, s)

	return nil
}
