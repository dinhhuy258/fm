package view

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/style"
	"github.com/dinhhuy258/gocui"
)

type SortAndFilterView struct {
	v            *View
	contentStyle style.TextStyle
}

func newSortAndFilterView(g *gocui.Gui, v *gocui.View) *SortAndFilterView {
	sfv := &SortAndFilterView{
		v:            newView(g, v),
		contentStyle: style.AttrBold,
	}

	sfv.v.SetTitle(" Sort & Filter ")
	sfv.SetSortAndFilter()

	return sfv
}

func (sfv *SortAndFilterView) SetSortAndFilter() {
	content := ""
	if !config.AppConfig.ShowHidden {
		content += "rel!^."
	}

	sfv.v.SetViewContent([]string{sfv.contentStyle.Sprint(content)})
}
