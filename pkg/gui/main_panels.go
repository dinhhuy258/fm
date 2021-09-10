package gui

import (
	"strconv"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
)

func (gui *Gui) RenderDir(dir *fs.Directory, focusIdx int) error {
	nodeSize := len(dir.Nodes)
	lines := make([]string, nodeSize)
	config := config.AppConfig

	for i, node := range dir.Nodes {
		fileIcon := config.FileIcon + " "
		if node.IsDir {
			fileIcon = config.FolderIcon + " "
		}

		var path string
		if i == focusIdx {
			path = config.FocusPrefix + fileIcon + node.RelativePath + config.FocusSuffix
		} else {
			path = "  " + fileIcon + node.RelativePath
		}

		if i == nodeSize-1 {
			path = config.TreeSuffix + path
		} else {
			path = config.TreePrefix + path
		}

		row := gui.MainRow.FileRow
		if node.IsDir {
			row = gui.MainRow.DirectoryRow
		}

		size := fs.Humanize(node.Size)
		index := strconv.Itoa(i-focusIdx) + "|" + strconv.Itoa(i)

		line, err := row.Sprint([]string{index, path, size})
		if err != nil {
			return err
		}

		lines[i] = line
	}

	gui.SetViewContent(gui.Views.Main, lines)

	return nil
}
