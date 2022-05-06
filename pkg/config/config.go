package config

import (
	"github.com/dinhhuy258/gocui"
	"github.com/gookit/color"
)

type CommandConfig struct {
	Name string
	Help string
	Args []interface{}
}

type KeyBindingsConfig struct {
	OnKeys map[string]*CommandConfig
}

type ModeConfig struct {
	Name        string
	KeyBindings KeyBindingsConfig
}

type Config struct {
	SelectionColor   color.Color
	DirectoryColor   color.Color
	SizeStyle        color.Color
	LogErrorColor    color.Color
	LogWarningColor  color.Color
	LogInfoColor     color.Color
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
	FolderIcon       string
	FileIcon         string
	LogErrorFormat   string
	LogWarningFormat string
	LogInfoFormat    string
	ModeConfigs      []ModeConfig
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
		ModeConfigs: []ModeConfig{
			{
				Name: "go-to",
				KeyBindings: KeyBindingsConfig{
					OnKeys: map[string]*CommandConfig{
						"~": {
							Name: "ChangeDirectory",
							Help: "Home",
							Args: []interface{}{"/Users/dinhhuy258"},
						},
						"d": {
							Name: "ChangeDirectory",
							Help: "Downloads",
							Args: []interface{}{"/Users/dinhhuy258/Downloads"},
						},
						"D": {
							Name: "ChangeDirectory",
							Help: "Documents",
							Args: []interface{}{"/Users/dinhhuy258/Documents"},
						},
						"w": {
							Name: "ChangeDirectory",
							Help: "Workspace",
							Args: []interface{}{"/Users/dinhhuy258/Workspace"},
						},
						"h": {
							Name: "ChangeDirectory",
							Help: "Desktop",
							Args: []interface{}{"/Users/dinhhuy258/Desktop"},
						},
						"q": {
							Name: "Quit",
							Help: "quit",
						},
						"esc": {
							Name: "PopMode",
							Help: "cancel",
						},
					},
				},
			},
		},
	}
}
