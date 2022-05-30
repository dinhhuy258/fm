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

	headerRow *style.Row
}

func newExplorerHeaderView(v *gocui.View) *ExplorerHeaderView {
	ehv := &ExplorerHeaderView{
		View:      newView(v),
		headerRow: newExplorerRow(),
	}

	return ehv
}

func (ehv *ExplorerHeaderView) layout() error {
	x, _ := ehv.Size()
	ehv.headerRow.SetWidth(x)

	rowString, err := ehv.headerRow.Sprint(
		[]string{
			config.AppConfig.IndexHeader,
			config.AppConfig.PathHeader,
			config.AppConfig.FileModeHeader,
			config.AppConfig.SizeHeader,
		},
	)
	if err != nil {
		return err
	}

	ehv.SetContent(rowString)

	return nil
}

type ExplorerView struct {
	*View

	explorerRow        *style.Row
	iconTextStyles     map[string]style.TextStyle
	fileTextStyle      style.TextStyle
	directoryTextStyle style.TextStyle
}

func newExplorerView(v *gocui.View) *ExplorerView {
	cfg := config.AppConfig

	ev := &ExplorerView{
		View:        newView(v),
		explorerRow: newExplorerRow(),
	}

	ev.Frame = false
	ev.Highlight = true

	ev.SelBgColor = style.StringToGoCuiColor(config.AppConfig.FocusBg)
	ev.SelFgColor = style.StringToGoCuiColor(config.AppConfig.FocusFg)

	ev.directoryTextStyle = style.FromBasicFg(style.StringToColor(cfg.NodeTypesConfig.Directory.Color))
	ev.fileTextStyle = style.FromBasicFg(style.StringToColor(cfg.NodeTypesConfig.File.Color))
	ev.iconTextStyles = map[string]style.TextStyle{}

	extensionNodeTypesConfig := cfg.NodeTypesConfig.Extensions

	ev.iconTextStyles[cfg.NodeTypesConfig.File.Icon] = ev.fileTextStyle
	ev.iconTextStyles[cfg.NodeTypesConfig.Directory.Icon] = ev.directoryTextStyle
	for _, ntc := range extensionNodeTypesConfig {
		if ntc.Color != "" {
			ev.iconTextStyles[ntc.Icon] = style.FromBasicFg(style.StringToColor(ntc.Color))
		} else {
			ev.iconTextStyles[ntc.Icon] = ev.fileTextStyle
		}
	}

	return ev
}

func newExplorerRow() *style.Row {
	r := &style.Row{}

	r.AddCell(config.AppConfig.IndexPercentage, true)
	r.AddCell(config.AppConfig.PathPercentage, true)
	r.AddCell(config.AppConfig.FileModePercentage, true)
	r.AddCell(config.AppConfig.SizePercentage, false)

	return r
}

func (ev *ExplorerView) UpdateView(entries []fs.IEntry, selections set.Set[string], focus int) {
	entriesSize := len(entries)
	lines := make([]string, entriesSize)
	cfg := config.AppConfig
	extensionNodeTypesConfig := cfg.NodeTypesConfig.Extensions
	var nameTextStyle style.TextStyle

	for idx, entry := range entries {
		var icon string
		if nodeTypeConfig, hasConfig := extensionNodeTypesConfig[entry.GetExt()]; hasConfig &&
			nodeTypeConfig.Icon != "" {
			icon = nodeTypeConfig.Icon
			nameTextStyle = ev.fileTextStyle
		} else if entry.IsDirectory() {
			icon = cfg.NodeTypesConfig.Directory.Icon
			nameTextStyle = ev.directoryTextStyle
		} else {
			nameTextStyle = ev.fileTextStyle
			icon = cfg.NodeTypesConfig.File.Icon
		}

		name := entry.GetName()
		var prefix string
		var suffix string
		var rowPrefix string

		switch {
		case idx == focus:
			prefix = cfg.FocusPrefix
			suffix = cfg.FocusSuffix
		case selections.Contains(entry.GetPath()):
			prefix = cfg.SelectionPrefix
			suffix = cfg.SelectionSuffix
		default:
			// TODO: Configurate these values
			prefix = "  "
			suffix = ""
		}

		if idx == entriesSize-1 {
			rowPrefix = cfg.PathSuffix
		} else {
			rowPrefix = cfg.PathPrefix
		}

		index := strconv.Itoa(idx + 1)
		fileMode := entry.GetFileMode()
		size := fs.Humanize(entry.GetSize())

		path := rowPrefix + ev.iconTextStyles[icon].Sprint(icon) + nameTextStyle.Sprint(prefix+name+suffix)
		line, err := ev.explorerRow.Sprint([]string{index, path, fileMode, size})
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
