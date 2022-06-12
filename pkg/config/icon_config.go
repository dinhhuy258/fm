package config

import icons "github.com/dinhhuy258/logo-ls/assets"

// getSpecialsNodeTypeConfig get default configuration for the special node type.
func getSpecialsNodeTypeConfig() map[string]*NodeTypeConfig {
	var specialsNodeTypeConfig map[string]*NodeTypeConfig = make(map[string]*NodeTypeConfig)
	for fileName, icon := range icons.Icon_FileName {
		specialsNodeTypeConfig[fileName] = &NodeTypeConfig{
			Icon: icon.GetGlyph(),
			Style: &StyleConfig{
				Fg: icon.GetHexColor(),
			},
		}
	}

	for dirName, icon := range icons.Icon_Dir {
		specialsNodeTypeConfig[dirName] = &NodeTypeConfig{
			Icon: icon.GetGlyph(),
			Style: &StyleConfig{
				Fg: icon.GetHexColor(),
			},
		}
	}

	return specialsNodeTypeConfig
}

// getExtensionsNodeTypeConfig get default configuration for the extensions node type.
func getExtensionsNodeTypeConfig() map[string]*NodeTypeConfig {
	var extensionsNodeTypeConfig map[string]*NodeTypeConfig = make(map[string]*NodeTypeConfig)
	for ext, icon := range icons.Icon_Ext {
		extensionsNodeTypeConfig[ext] = &NodeTypeConfig{
			Icon: icon.GetGlyph(),
			Style: &StyleConfig{
				Fg: icon.GetHexColor(),
			},
		}
	}

	return extensionsNodeTypeConfig
}
