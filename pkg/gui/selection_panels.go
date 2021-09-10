package gui

func (gui *Gui) RenderSelections(selections map[string]struct{}) error {
	s := make([]string, 0, len(selections))
	for k := range selections {
		s = append(s, k)
	}

	gui.SetViewContent(gui.Views.Selection, s)

	return nil
}
