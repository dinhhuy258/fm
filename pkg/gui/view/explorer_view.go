package view

import (
	"log"
	"strconv"
	"strings"

	set "github.com/deckarep/golang-set/v2"
	"github.com/dinhhuy258/fm/pkg/config"
	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view/style"
	"github.com/dinhhuy258/gocui"
)

type ExplorerHeaderView struct {
	*View

	headerRow *Row
}

func newExplorerHeaderView(v *gocui.View) *ExplorerHeaderView {
	ehv := &ExplorerHeaderView{
		View:      newView(v),
		headerRow: newExplorerRow(),
	}

	return ehv
}

func (ehv *ExplorerHeaderView) layout() error {
	explorerTableConfig := config.AppConfig.General.ExplorerTable

	indexHeaderTextStyle := style.FromStyleConfig(explorerTableConfig.IndexHeader.Style)
	nameHeaderTextStyle := style.FromStyleConfig(explorerTableConfig.NameHeader.Style)
	permissionsHeaderTextStyle := style.FromStyleConfig(explorerTableConfig.PermissionsHeader.Style)
	sizeHeaderTextStyle := style.FromStyleConfig(explorerTableConfig.SizeHeader.Style)

	x, _ := ehv.Size()
	ehv.headerRow.SetWidth(x)

	rowString, err := ehv.headerRow.Sprint(
		[]ColumnValue{
			[]ColumnValueComponent{
				{
					Value: explorerTableConfig.IndexHeader.Name,
					Style: &indexHeaderTextStyle,
				},
			},
			[]ColumnValueComponent{
				{
					Value: explorerTableConfig.NameHeader.Name,
					Style: &nameHeaderTextStyle,
				},
			},
			[]ColumnValueComponent{
				{
					Value: explorerTableConfig.PermissionsHeader.Name,
					Style: &permissionsHeaderTextStyle,
				},
			},
			[]ColumnValueComponent{
				{
					Value: explorerTableConfig.SizeHeader.Name,
					Style: &sizeHeaderTextStyle,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	ehv.SetContent(rowString)

	return nil
}

type nodeType struct {
	icon  string
	style style.TextStyle
}

type nodeTypes struct {
	file             nodeType
	directory        nodeType
	fileSymlink      nodeType
	directorySymlink nodeType
	extensions       map[string]nodeType
}

type ExplorerView struct {
	*View

	explorerRow *Row
	icons       nodeTypes

	defaultFileTextStyle      style.TextStyle
	defaultDirectoryTextStyle style.TextStyle

	focusTextStyle          style.TextStyle
	selectionTextStyle      style.TextStyle
	focusSelectionTextStyle style.TextStyle
}

func newExplorerView(v *gocui.View) *ExplorerView {
	nodeTypesConfig := config.AppConfig.NodeTypesConfig
	explorerTableConfig := config.AppConfig.General.ExplorerTable

	ev := &ExplorerView{
		View:        newView(v),
		explorerRow: newExplorerRow(),
	}

	ev.Frame = false

	ev.defaultFileTextStyle = style.FromStyleConfig(explorerTableConfig.DefaultUI.FileStyle)
	ev.defaultDirectoryTextStyle = style.FromStyleConfig(explorerTableConfig.DefaultUI.DirectoryStyle)

	ev.focusTextStyle = style.FromStyleConfig(explorerTableConfig.FocusUI.Style)
	ev.focusSelectionTextStyle = style.FromStyleConfig(explorerTableConfig.FocusSelectionUI.Style)
	ev.selectionTextStyle = style.FromStyleConfig(explorerTableConfig.SelectionUI.Style)

	ev.icons = nodeTypes{
		file: nodeType{
			icon:  nodeTypesConfig.File.Icon,
			style: style.FromStyleConfig(nodeTypesConfig.File.Style),
		},
		directory: nodeType{
			icon:  nodeTypesConfig.Directory.Icon,
			style: style.FromStyleConfig(nodeTypesConfig.Directory.Style),
		},
		fileSymlink: nodeType{
			icon:  nodeTypesConfig.FileSymlink.Icon,
			style: style.FromStyleConfig(nodeTypesConfig.FileSymlink.Style),
		},
		directorySymlink: nodeType{
			icon:  nodeTypesConfig.DirectorySymlink.Icon,
			style: style.FromStyleConfig(nodeTypesConfig.DirectorySymlink.Style),
		},
		extensions: map[string]nodeType{},
	}

	for ext, ntc := range nodeTypesConfig.Extensions {
		ev.icons.extensions[ext] = nodeType{
			icon:  ntc.Icon,
			style: style.FromStyleConfig(ntc.Style),
		}
	}

	return ev
}

func newExplorerRow() *Row {
	explorerTableConfig := config.AppConfig.General.ExplorerTable

	r := &Row{}

	r.AddColumn(explorerTableConfig.IndexHeader.Percentage, true)
	r.AddColumn(explorerTableConfig.NameHeader.Percentage, true)
	r.AddColumn(explorerTableConfig.PermissionsHeader.Percentage, true)
	r.AddColumn(explorerTableConfig.SizeHeader.Percentage, false)

	return r
}

func (ev *ExplorerView) UpdateView(entries []fs.IEntry, selections set.Set[string], focus int) {
	explorerTableConfig := config.AppConfig.General.ExplorerTable

	entriesSize := len(entries)
	lines := make([]string, entriesSize)

	for idx, entry := range entries {
		isEntryFocused := idx == focus
		isEntrySelected := selections.Contains(entry.GetPath())

		entryIcon := ev.getEntryIcon(entry, isEntryFocused, isEntrySelected)

		name := entry.GetName()

		var prefix, suffix, entryTreePrefix string

		var entryTextStyle style.TextStyle

		switch {
		case isEntryFocused && isEntrySelected:
			prefix = explorerTableConfig.FocusSelectionUI.Prefix
			suffix = explorerTableConfig.FocusSelectionUI.Suffix

			entryTextStyle = ev.focusSelectionTextStyle
		case isEntryFocused:
			prefix = explorerTableConfig.FocusUI.Prefix
			suffix = explorerTableConfig.FocusUI.Suffix

			entryTextStyle = ev.focusTextStyle
		case isEntrySelected:
			prefix = explorerTableConfig.SelectionUI.Prefix
			suffix = explorerTableConfig.SelectionUI.Suffix

			entryTextStyle = ev.selectionTextStyle
		default:
			prefix = explorerTableConfig.DefaultUI.Prefix
			suffix = explorerTableConfig.DefaultUI.Suffix

			if entry.IsDirectory() {
				entryTextStyle = ev.defaultDirectoryTextStyle
			} else {
				entryTextStyle = ev.defaultFileTextStyle
			}
		}

		switch {
		case idx == entriesSize-1:
			entryTreePrefix = explorerTableConfig.LastEntryPrefix
		case idx == 0:
			entryTreePrefix = explorerTableConfig.FirstEntryPrefix
		default:
			entryTreePrefix = explorerTableConfig.EntryPrefix
		}

		index := strconv.Itoa(idx + 1)
		permissions := entry.GetPermissions()
		size := fs.Humanize(entry.GetSize())

		line, err := ev.explorerRow.Sprint([]ColumnValue{
			index,
			[]ColumnValueComponent{
				{
					Value: entryTreePrefix,
					Style: nil,
				},
				{
					Value: prefix,
					Style: &entryTextStyle,
				},
				{
					Value: entryIcon.icon,
					Style: &entryIcon.style,
				},
				{
					Value: " ",
					Style: nil,
				},
				{
					Value: name,
					Style: &entryTextStyle,
				},
				{
					Value: suffix,
					Style: &entryTextStyle,
				},
			},
			permissions,
			size,
		})
		if err != nil {
			log.Fatalf("failed to sprint row %v", err)
		}

		lines[idx] = line
	}

	ev.SetContent(strings.Join(lines, "\n"))
}

func (ev *ExplorerView) layout() error {
	x, _ := ev.Size()
	ev.explorerRow.SetWidth(x)

	return nil
}

func (ev *ExplorerView) getEntryIcon(entry fs.IEntry, isEntryFocused, isEntrySelected bool) nodeType {
	var icon nodeType

	extensionIcon, hasExtIcon := ev.icons.extensions[entry.GetExt()]
	fileIcon, hasFileIcon := ev.icons.extensions[entry.GetName()]

	switch {
	case entry.IsSymlink() && entry.IsDirectory():
		icon = ev.icons.directorySymlink
	case entry.IsSymlink():
		icon = ev.icons.fileSymlink
	case entry.IsDirectory():
		icon = ev.icons.directory
	case hasExtIcon:
		icon = extensionIcon
	case hasFileIcon:
		icon = fileIcon
	default:
		icon = ev.icons.file
	}

	switch {
	case isEntrySelected && isEntryFocused:
		icon.style = ev.focusSelectionTextStyle
	case isEntrySelected:
		icon.style = ev.selectionTextStyle
	case isEntryFocused:
		icon.style = ev.focusTextStyle
	}

	return icon
}
