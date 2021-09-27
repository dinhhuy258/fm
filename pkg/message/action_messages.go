package message

import (
	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func NewFile(ctx ctx.Context, _ ...interface{}) error {
	ctx.Gui().Views.Input.SetInput("new file")

	go func() {
		ans := ctx.Gui().Views.Input.GetAnswer()

		ctx.Gui().Views.Main.SetAsCurrentView()

		if ans != "" {
			ctx.Gui().Views.Log.SetLog(ans, view.LogLevel(view.ERROR))
		} else {
			ctx.Gui().Views.Log.SetViewOnTop()
		}
	}()

	return nil
}
