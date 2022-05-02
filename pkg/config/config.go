package config

import (
	"github.com/dinhhuy258/gocui"
	"github.com/gookit/color"
)

type Config struct {
	ShowHidden       bool
	IndexHeader      string
	IndexPercentage  int
	PathHeader       string
	PathPercentage   int
	SizeHeader       string
	SizePercentage   int
	PathPrefix       string
	PathSuffix       string
	FocusPrefix      string
	FocusSuffix      string
	FocusBg          gocui.Attribute
	FocusFg          gocui.Attribute
	SelectionPrefix  string
	SelectionSuffix  string
	SelectionColor   color.Color
	FolderIcon       string
	FileIcon         string
	DirectoryColor   color.Color
	SizeStyle        color.Color
	LogErrorFormat   string
	LogErrorColor    color.Color
	LogWarningFormat string
	LogWarningColor  color.Color
	LogInfoFormat    string
	LogInfoColor     color.Color
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		ShowHidden:       false,
		IndexHeader:      "index",
		IndexPercentage:  10,
		PathHeader:       "╭──── path",
		PathPercentage:   70,
		SizeHeader:       "size",
		SizePercentage:   20,
		PathPrefix:       "├─",
		PathSuffix:       "╰─",
		FocusPrefix:      "▸[",
		FocusSuffix:      "]",
		FocusBg:          gocui.ColorDefault,
		FocusFg:          gocui.ColorBlue,
		SelectionPrefix:  "{",
		SelectionSuffix:  "}",
		SelectionColor:   color.Green,
		FolderIcon:       "",
		FileIcon:         "",
		DirectoryColor:   color.Cyan,
		SizeStyle:        color.White,
		LogErrorFormat:   "[ERROR] ",
		LogErrorColor:    color.Red,
		LogWarningFormat: "[WARNING] ",
		LogWarningColor:  color.Yellow,
		LogInfoFormat:    "[INFO] ",
		LogInfoColor:     color.Green,
	}
}
