package message

import (
	"fmt"
	"github.com/dinhhuy258/fm/pkg/app/context"
	"path"
	"strings"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func NewFile(ctx context.Context, _ ...interface{}) error {
	gui.GetGui().Views.Input.SetInput("new file")

	go func() {
		name := gui.GetGui().Views.Input.GetAnswer()

		gui.GetGui().Views.Main.SetAsCurrentView()

		if name == "" {
			gui.GetGui().Views.Log.SetViewOnTop()

			return
		}

		var err error

		if strings.HasSuffix(name, "/") {
			err = fs.CreateDirectory(name)
		} else {
			err = fs.CreateFile(name)
		}

		if err != nil {
			gui.GetGui().Views.Log.SetLog(fmt.Sprintf("Failed to create file %s", name), view.LogLevel(view.ERROR))
		} else {
			gui.GetGui().Views.Log.SetLog(fmt.Sprintf("File %s were created successfully", name),
				view.LogLevel(view.INFO))
			_ = Refresh(ctx, path.Join(fs.GetFileManager().Dir.Path, name))
		}
	}()

	return nil
}
