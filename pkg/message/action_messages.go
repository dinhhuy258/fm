package message

import (
	"fmt"
	"strings"

	"github.com/dinhhuy258/fm/pkg/ctx"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func NewFile(ctx ctx.Context, _ ...interface{}) error {
	ctx.Gui().Views.Input.SetInput("new file")

	go func() {
		name := ctx.Gui().Views.Input.GetAnswer()

		ctx.Gui().Views.Main.SetAsCurrentView()

		if name == "" {
			ctx.Gui().Views.Log.SetViewOnTop()

			return
		}

		var err error

		if strings.HasSuffix(name, "/") {
			err = fs.CreateDirectory(name)
		} else {
			err = fs.CreateFile(name)
		}

		if err != nil {
			ctx.Gui().Views.Log.SetLog(fmt.Sprintf("Failed to create file %s", name), view.LogLevel(view.ERROR))
		} else {
			ctx.Gui().Views.Log.SetLog(fmt.Sprintf("File %s were created successfully", name),
				view.LogLevel(view.INFO))
			_ = Refresh(ctx)
		}
	}()

	return nil
}
