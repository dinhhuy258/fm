local fm = fm

fm.general.logInfoUi = {
	prefix = "[Info] ",
	suffix = "",
	style = {
		fg = "green",
	},
}

fm.general.logWarningUi = {
	prefix = "[Warning] ",
	suffix = "",
	style = {
		fg = "yellow",
	},
}

fm.general.logErrorUi = {
	prefix = "[Error] ",
	suffix = "",
	style = {
		fg = "red",
	},
}

fm.general.explorerTable.defaultUi = {
	prefix = "  ",
	suffix = "",
	file_style = {
		fg = "white",
	},
	directory_style = {
		fg = "cyan",
	},
}

fm.general.explorerTable.focusUi = {
	prefix = "▸[",
	suffix = "]",
	style = {
		fg = "white",
		decorations = {
			"bold",
		},
	},
}

fm.general.explorerTable.selectionUi = {
	prefix = " {",
	suffix = "}",
	style = {
		fg = "green",
	},
}

fm.general.explorerTable.focusSelectionUi = {
	prefix = "▸[",
	suffix = "]",
	style = {
		fg = "green",
		decorations = {
			"bold",
		},
	},
}

fm.general.explorerTable.indexHeader = {
	name = "index",
	percentage = 10,
}

fm.general.explorerTable.nameHeader = {
	name = "┌──── name",
	percentage = 65,
}

fm.general.explorerTable.permissionsHeader = {
	name = "permissions",
	percentage = 15,
}

fm.general.explorerTable.sizeHeader = {
	name = "size",
	percentage = 10,
}
