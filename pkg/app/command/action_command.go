package command

import (
	"fmt"
	"path"
	"strings"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func NewFile(app IApp, _ ...interface{}) error {
	gui.GetGui().SetInput("new file", func(name string) {
		if name == "" {
			gui.GetGui().SetLog("File name is empty", view.LogLevel(view.WARNING))

			return
		}

		var err error

		if strings.HasSuffix(name, "/") {
			err = fs.CreateDirectory(name)
		} else {
			err = fs.CreateFile(name)
		}

		if err != nil {
			gui.GetGui().SetLog(fmt.Sprintf("Failed to create file %s", name), view.LogLevel(view.ERROR))
		} else {
			gui.GetGui().SetLog(fmt.Sprintf("File %s were created successfully", name),
				view.LogLevel(view.INFO))
			_ = Refresh(app, path.Join(fs.GetFileManager().Dir.Path, name))
		}
	})

	return nil
}
