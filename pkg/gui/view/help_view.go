package view

import (
	"fmt"
	"log"

	"github.com/dinhhuy258/fm/pkg/gui/view/row"
	"github.com/dinhhuy258/gocui"
)

type HelpView struct {
	*View

	helpRow *row.Row
}

func newHelpView(g *gocui.Gui, v *gocui.View) *HelpView {
	hv := &HelpView{
		View: newView(g, v),
	}

	hv.helpRow = &row.Row{}
	hv.helpRow.AddCell(35, true, nil)
	hv.helpRow.AddCell(65, true, nil)

	return hv
}

func (hv *HelpView) layout() {
	x, _ := hv.v.Size()
	hv.helpRow.SetWidth(x)
}

func (hv *HelpView) UpdateView(title string, keys []string, msgs []string) {
	lines := make([]string, 0, len(keys))

	for i := 0; i < len(keys); i++ {
		key := keys[i]
		msg := msgs[i]

		line, err := hv.helpRow.Sprint([]string{key, msg})
		if err != nil {
			log.Fatalf("failed to set content for help view %v", err)
		}

		lines = append(lines, line)
	}

	hv.SetViewContent(lines)
	hv.v.Title = fmt.Sprintf(" Help [%s] ", title)
}
