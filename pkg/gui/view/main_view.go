package view

import (
	"strconv"

	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/row"
	"github.com/dinhhuy258/gocui"
)

type MainView struct {
	v       *View
	mainRow *row.MainRow
}

func newMainView(g *gocui.Gui, v *gocui.View) *MainView {
	mv := &MainView{
		v:       newView(g, v),
		mainRow: row.NewMainRow(),
	}

	mv.v.v.Frame = false
	mv.v.v.Highlight = true
	mv.v.v.SelBgColor = config.AppConfig.FocusBg
	mv.v.v.SelFgColor = config.AppConfig.FocusFg

	return mv
}

func (mv *MainView) OnResize() {
	x, _ := mv.v.v.Size()
	mv.mainRow.SetWidth(x)
}

func (mv *MainView) RenderDir(dir *fs.Directory, selections map[string]struct{}, focusIdx int) error {
	nodeSize := len(dir.Nodes)
	lines := make([]string, nodeSize)
	config := config.AppConfig

	for i, node := range dir.Nodes {
		fileIcon := config.FileIcon + " "
		if node.IsDir {
			fileIcon = config.FolderIcon + " "
		}

		_, isSelected := selections[node.AbsolutePath]

		var path string

		switch {
		case i == focusIdx:
			path = config.FocusPrefix + fileIcon + node.RelativePath + config.FocusSuffix
		case isSelected:
			path = config.SelectionPrefix + fileIcon + node.RelativePath + config.SelectionSuffix
		default:
			path = "  " + fileIcon + node.RelativePath
		}

		if i == nodeSize-1 {
			path = config.PathSuffix + path
		} else {
			path = config.PathPrefix + path
		}

		row := mv.mainRow.FileRow
		if isSelected {
			row = mv.mainRow.SelectionRow
		} else if node.IsDir {
			row = mv.mainRow.DirectoryRow
		}

		size := fs.Humanize(node.Size)
		index := strconv.Itoa(i-focusIdx) + "|" + strconv.Itoa(i)

		line, err := row.Sprint([]string{index, path, size})
		if err != nil {
			return err
		}

		lines[i] = line
	}

	mv.v.SetViewContent(lines)

	return nil
}

func (mv *MainView) SetOrigin(x, y int) error {
	return mv.v.SetOrigin(x, y)
}

func (mv *MainView) SetCursor(x, y int) error {
	return mv.v.SetCursor(x, y)
}

func (mv *MainView) NextCursor() error {
	return mv.v.NextCursor()
}

func (mv *MainView) PreviousCursor() error {
	return mv.v.PreviousCursor()
}
