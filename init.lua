local fm = fm

fm.general.log_info_ui = {
	prefix = "[Info] ",
	suffix = "",
	style = {
		fg = "green",
	},
}

fm.general.log_warning_ui = {
	prefix = "[Warning] ",
	suffix = "",
	style = {
		fg = "yellow",
	},
}

fm.general.log_error_ui = {
	prefix = "[Error] ",
	suffix = "",
	style = {
		fg = "red",
	},
}

fm.general.explorer.default_ui = {
	prefix = "  ",
	suffix = "",
	file_style = {
		fg = "white",
	},
	directory_style = {
		fg = "cyan",
	},
}

fm.general.explorer.focus_ui = {
	prefix = "▸[",
	suffix = "]",
	style = {
		fg = "white",
		decorations = {
			"bold",
		},
	},
}

fm.general.explorer.selection_ui = {
	prefix = " {",
	suffix = "}",
	style = {
		fg = "green",
	},
}

fm.general.explorer.focus_selection_ui = {
	prefix = "▸[",
	suffix = "]",
	style = {
		fg = "green",
		decorations = {
			"bold",
		},
	},
}

fm.general.explorer.index_header = {
	name = "index",
	percentage = 10,
}

fm.general.explorer.name_header = {
	name = "┌──── name",
	percentage = 65,
}

fm.general.explorer.permissions_header = {
	name = "permissions",
	percentage = 15,
}

fm.general.explorer.size_header = {
	name = "size",
	percentage = 10,
}
