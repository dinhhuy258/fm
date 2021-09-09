package gui

import (
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
)

func (gui *Gui) RenderDir(dir *fs.Directory, selectedIdx int) {
	nodeSize := len(dir.Nodes)
	lines := make([]string, nodeSize)
	config := config.AppConfig

	for i, node := range dir.Nodes {
		fileIcon := config.FileIcon + " "
		if node.IsDir {
			fileIcon = config.FolderIcon + " "
		}

		if i == selectedIdx {
			lines[i] = config.FocusPrefix + fileIcon + node.RelativePath + config.FocusSuffix
		} else {
			lines[i] = "  " + fileIcon + node.RelativePath
		}

		if i == nodeSize-1 {
			lines[i] = config.TreeSuffix + lines[i]
		} else {
			lines[i] = config.TreePrefix + lines[i]
		}

		if node.IsDir {
			lines[i] = config.DirectoryStyle.Sprint(lines[i])
		}
	}

	gui.SetViewContent(gui.Views.Main, lines)
}
