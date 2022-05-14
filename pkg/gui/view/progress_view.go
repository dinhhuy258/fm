package view

import (
	"fmt"

	"github.com/dinhhuy258/gocui"
)

const (
	progressEmpty string = "░"
	progressFull  string = "█"
)

type ProgressView struct {
	v       *View
}

func newProgressView(g *gocui.Gui, v *gocui.View) *ProgressView {
	pv := &ProgressView{
		v: newView(g, v),
	}

	return pv
}

func (pv *ProgressView) SetViewOnTop() {
	pv.v.SetViewOnTop()
}

func (pv *ProgressView) UpdateView(current int, total int) {
	percent := float32(current) / float32(total)

	x, _ := pv.v.Size()

	progressBar := ""
	fullCount := int(float32(x) * percent)
	emptyCount := x - fullCount

	for i := 0; i < fullCount; i++ {
		progressBar += progressFull
	}

	for i := 0; i < emptyCount; i++ {
		progressBar += progressEmpty
	}

	pv.v.v.Title = fmt.Sprintf(" Progress (%s) ", fmt.Sprintf("%0.0f%%", percent*100))
	pv.v.SetViewContent([]string{progressBar})
}
