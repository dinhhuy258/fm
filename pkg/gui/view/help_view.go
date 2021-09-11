package view

import (
	"github.com/dinhhuy258/fm/pkg/row"
	"github.com/dinhhuy258/gocui"
)

type HelpView struct {
	v       *View
	helpRow *row.Row
}

func newHelpView(g *gocui.Gui, v *gocui.View) *HelpView {
	hv := &HelpView{
		v: newView(g, v),
	}

	hv.helpRow = &row.Row{}
	hv.helpRow.AddCell(35, true, nil)
	hv.helpRow.AddCell(65, true, nil)

	hv.v.v.Title = " Help "

	return hv
}

func (hv *HelpView) layout() {
	x, _ := hv.v.v.Size()
	hv.helpRow.SetWidth(x)
}

func (hv *HelpView) SetViewContent(keys []string, helps []string) error {
	lines := make([]string, 0, len(keys))

	for i := 0; i < len(keys); i++ {
		key := keys[i]
		help := helps[i]

		line, err := hv.helpRow.Sprint([]string{key, help})
		if err != nil {
			return err
		}

		lines = append(lines, line)
	}

	hv.v.SetViewContent(lines)

	return nil
}
