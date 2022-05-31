package view

import (
	"fmt"
	"log"
	"strings"

	"github.com/dinhhuy258/gocui"
)

type HelpView struct {
	*View

	helpRow *Row
}

func newHelpView(v *gocui.View) *HelpView {
	hv := &HelpView{
		View: newView(v),
	}

	hv.helpRow = &Row{}
	hv.helpRow.AddColumn(35, true)
	hv.helpRow.AddColumn(65, true)

	return hv
}

func (hv *HelpView) UpdateView(title string, helpKeys []string, helpMsgs []string) {
	lines := make([]string, 0, len(helpKeys))

	for i := 0; i < len(helpKeys); i++ {
		helpKey := helpKeys[i]
		helpMsg := helpMsgs[i]

		line, err := hv.helpRow.Sprint([]ColumnValue{helpKey, helpMsg})
		if err != nil {
			log.Fatalf("failed to set content for help view %v", err)
		}

		lines = append(lines, line)
	}

	hv.SetContent(strings.Join(lines, "\n"))
	hv.Title = fmt.Sprintf(" Help [%s] ", title)
}

func (hv *HelpView) layout() error {
	x, _ := hv.Size()
	hv.helpRow.SetWidth(x)

	return nil
}
