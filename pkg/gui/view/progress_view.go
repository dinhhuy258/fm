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
	total   int
	current int
}

func newProgressView(g *gocui.Gui, v *gocui.View) *ProgressView {
	pv := &ProgressView{
		v: newView(g, v),
	}

	return pv
}

func (pv *ProgressView) StartProgress(total int) error {
	pv.total = total
	pv.AddCurrent(0)

	return pv.v.SetViewOnTop()
}

func (pv *ProgressView) AddCurrent(current int) {
	pv.current += current
	percent := float32(current) / float32(pv.total)

	x, _ := pv.v.Size()
	pv.v.v.Title = fmt.Sprintf(" Progress (%s) ", fmt.Sprintf("%0.0f%%", percent*100))
	progressBar := ""

	fullCount := int(float32(x) * percent)
	emptyCount := x - fullCount

	for i := 0; i < fullCount; i++ {
		progressBar += progressFull
	}

	for i := 0; i < emptyCount; i++ {
		progressBar += progressEmpty
	}

	pv.v.SetViewContent([]string{progressBar})
}

func (pv *ProgressView) IsFinished() bool {
	return pv.total == pv.current
}
