package view

import (
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/gui/view/row"
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

	return hv
}

func (hv *HelpView) layout() {
	x, _ := hv.v.v.Size()
	hv.helpRow.SetWidth(x)
}

func (hv *HelpView) SetTitle(mode string) {
	hv.v.v.Title = fmt.Sprintf(" Help [%s] ", mode)
}

func (hv *HelpView) SetHelp(keys []string, helps []string) {
	lines := make([]string, 0, len(keys))

	for i := 0; i < len(keys); i++ {
		key := keys[i]
		help := helps[i]

		line, err := hv.helpRow.Sprint([]string{key, help})
		if err != nil {
			log.Fatalf("failed to set content for help view %v", err)
		}

		lines = append(lines, line)
	}

	hv.v.SetViewContent(lines)
}
