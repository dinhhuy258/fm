package controller

import "github.com/dinhhuy258/fm/pkg/gui/view"

type ProgressController struct {
	*BaseController

	total   int
	current int

	view *view.ProgressView
}

func newProgressController(baseController *BaseController,
	view *view.ProgressView) *ProgressController {
	return &ProgressController{
		BaseController: baseController,
		view:           view,

		total:   0,
		current: 0,
	}
}

func (pc *ProgressController) SetView(view *view.ProgressView) {
	pc.view = view
}

func (pc *ProgressController) StartProgress(total int) {
	pc.total = total
	pc.current = 0

	pc.view.UpdateView(pc.current, pc.total)
	pc.view.SetViewOnTop()
}

func (pc *ProgressController) UpdateProgress() {
	pc.addCurrent(1)
}

func (pc *ProgressController) IsProgressFinished() bool {
	return pc.total == pc.current
}

func (pc *ProgressController) addCurrent(current int) {
	pc.current += current

	pc.view.UpdateView(pc.current, pc.total)
}
