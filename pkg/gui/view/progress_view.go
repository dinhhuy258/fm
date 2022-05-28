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
	*View
}

func newProgressView(v *gocui.View) *ProgressView {
	pv := &ProgressView{
		newView(v),
	}

	return pv
}

func (pv *ProgressView) UpdateView(current int, total int) {
	percent := float32(current) / float32(total)

	x, _ := pv.Size()

	progressBar := ""
	fullCount := int(float32(x) * percent)
	emptyCount := x - fullCount

	for i := 0; i < fullCount; i++ {
		progressBar += progressFull
	}

	for i := 0; i < emptyCount; i++ {
		progressBar += progressEmpty
	}

	pv.Title = fmt.Sprintf(" Progress (%s) ", fmt.Sprintf("%0.0f%%", percent*100))
	pv.SetViewContent([]string{progressBar})
}
