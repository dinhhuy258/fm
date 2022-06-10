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
		fg = "blue",
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

fm.general.explorerTable.firstEntryPrefix = "├─"
fm.general.explorerTable.entryPrefix = "├─"
fm.general.explorerTable.lastEntryPrefix = "└─"

fm.general.sorting = {
	sortType = "dirFirst",
	reverse = false,
	ignoreCase = true,
	ignoreDiacritics = true,
}

fm.modes.customs["go-to"] = {
	name = "go-to",
	keyBindings = {
		default = {
			help = "cancel",
			messages = {
				{
					name = "SwitchMode",
					args = {
						"default",
					},
				},
			},
		},
		onKeys = {
			["~"] = {
				help = "Home",
				messages = {
					{
						name = "ChangeDirectory",
						args = {
							"/Users/dinhhuy258",
						},
					},
					{
						name = "SwitchMode",
						args = {
							"default",
						},
					},
				},
			},
			w = {
				help = "Workspace",
				messages = {
					{
						name = "ChangeDirectory",
						args = {
							"/Users/dinhhuy258/Workspace",
						},
					},
					{
						name = "SwitchMode",
						args = {
							"default",
						},
					},
				},
			},
			D = {
				help = "Documents",
				messages = {
					{
						name = "ChangeDirectory",
						args = {
							"/Users/dinhhuy258/Documents",
						},
					},
					{
						name = "SwitchMode",
						args = {
							"default",
						},
					},
				},
			},
			d = {
				help = "Downloads",
				messages = {
					{
						name = "ChangeDirectory",
						args = {
							"/Users/dinhhuy258/Downloads",
						},
					},
					{
						name = "SwitchMode",
						args = {
							"default",
						},
					},
				},
			},
			g = {
				help = "focus first",
				messages = {
					{
						name = "FocusFirst",
					},
					{
						name = "SwitchMode",
						args = {
							"default",
						},
					},
				},
			},
			h = {
				help = "Desktop",
				messages = {
					{
						name = "ChangeDirectory",
						args = {
							"/Users/dinhhuy258/Desktop",
						},
					},
					{
						name = "SwitchMode",
						args = {
							"default",
						},
					},
				},
			},
			["ctrl+c"] = {
				help = "quit",
				messages = {
					{
						name = "Quit",
					},
				},
			},
		},
	},
}

fm.modes.customs.yarn = {
	name = "yarn",
	keyBindings = {
		default = {
			help = "cancel",
			messages = {
				{
					name = "SwitchMode",
					args = {
						"default",
					},
				},
			},
		},
		onKeys = {
			p = {
				help = "yarn path",
				messages = {
					{
						name = "BashExecSilently",
						args = {
							[===[
              focus_path="${FM_FOCUS_PATH}"
              if [ "${focus_path}" ]; then
                echo -n ${focus_path} | pbcopy -selection clipboard
                echo LogSuccess "'"${focus_path} was coppied to clipboard"'" >> "${FM_PIPE_MSG_IN:?}"
              fi

              echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
              ]===],
						},
					},
				},
			},
			n = {
				help = "yarn name",
				messages = {
					{
						name = "BashExecSilently",
						args = {
							[===[
              focus_path="${FM_FOCUS_PATH}"
              if [ "${focus_path}" ]; then
                echo -n $(basename ${focus_path}) | pbcopy -selection clipboard
                echo LogSuccess "'"$(basename "${focus_path}") was coppied to clipboard"'" >> "${FM_PIPE_MSG_IN:?}"
              fi

              echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
              ]===],
						},
					},
				},
			},
			d = {
				help = "yarn directory",
				messages = {
					{
						name = "BashExecSilently",
						args = {
							[===[
              focus_path="${FM_FOCUS_PATH}"
              if [ "${focus_path}" ]; then
                echo -n $(dirname ${focus_path}) | pbcopy -selection clipboard
                echo LogSuccess "'"$(dirname "${focus_path}") was coppied to clipboard"'" >> "${FM_PIPE_MSG_IN:?}"
              fi

              echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
              ]===],
						},
					},
				},
			},
			["ctrl+c"] = {
				help = "quit",
				messages = {
					{
						name = "Quit",
					},
				},
			},
		},
	},
}

fm.modes.customs["mark-save"] = {
	name = "mark-save",
	keyBindings = {
		default = {
			help = "save",
			messages = {
				{
					name = "SetInputBuffer",
					args = {
						"",
					},
				},
				{
					name = "UpdateInputBufferFromKey",
				},
				{
					name = "BashExecSilently",
					args = {
						[===[
            session_path="${FM_SESSION_PATH:?}"
            mark_file="${session_path}/mark"

            focus_path="${FM_FOCUS_PATH:?}"
            key="${FM_INPUT_BUFFER:?}"

            # create a mark file if not exists
            touch ${mark_file}
            # remove conflict mark key in the mark file
            marks=$(sed "/^"${key}";/d" < "${mark_file:?}")
            printf "%s\n" "${marks[@]}" > ${mark_file}
            # add new mark key to the mark file
            echo "${key};${focus_path}" >> ${mark_file}

            echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
            ]===],
					},
				},
			},
		},
		onKeys = {
			esc = {
				help = "cancel",
				messages = {
					{
						name = "SwitchMode",
						args = {
							"default",
						},
					},
				},
			},
			["ctrl+c"] = {
				help = "quit",
				messages = {
					{
						name = "Quit",
					},
				},
			},
		},
	},
}

fm.modes.customs["mark-load"] = {
	name = "mark-load",
	keyBindings = {
		default = {
			help = "load",
			messages = {
				{
					name = "SetInputBuffer",
					args = {
						"",
					},
				},
				{
					name = "UpdateInputBufferFromKey",
				},
				{
					name = "BashExecSilently",
					args = {
						[===[
            session_path="${FM_SESSION_PATH:?}"
            mark_file="${session_path}/mark"

            # create a mark file if not exists
            touch ${mark_file}

            pressed_key="${FM_INPUT_BUFFER}"

            # go through the mark file and read the value to variable marks
            (while IFS= read -r line; do
            key=$(echo ${line:?} | cut -d ";" -f 1)
            path=$(echo ${line:?} | cut -d ";" -f 2)

            if [ "${key}" = "${pressed_key}" ]; then
              echo FocusPath "'""${path}""'" >> "${FM_PIPE_MSG_IN:?}"
              break
            fi
            done < "${mark_file:?}")

            echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
            ]===],
					},
				},
			},
		},
		onKeys = {
			esc = {
				help = "cancel",
				messages = {
					{
						name = "SwitchMode",
						args = {
							"default",
						},
					},
				},
			},
			["ctrl+c"] = {
				help = "quit",
				messages = {
					{
						name = "Quit",
					},
				},
			},
		},
	},
}

fm.modes.customs.open = {
	name = "yarn",
	keyBindings = {
		default = {
			help = "cancel",
			messages = {
				{
					name = "SwitchMode",
					args = {
						"default",
					},
				},
			},
		},
		onKeys = {
			f = {
				help = "open in Finder",
				messages = {
					{
						name = "BashExecSilently",
						args = {
							[===[
              focus_path="${FM_FOCUS_PATH}"
              if [ "${focus_path}" ]; then
                open -R "${focus_path}"
              fi

              echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
              ]===],
						},
					},
				},
			},
			c = {
				help = "open in Visual Studio Code",
				messages = {
					{
						name = "BashExecSilently",
						args = {
							[===[
              focus_path="${FM_FOCUS_PATH}"
              if [ "${focus_path}" ]; then
                open -a "Visual Studio Code" "${focus_path}"
              fi

              echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
              ]===],
						},
					},
				},
			},
			v = {
				help = "open in NeoVim",
				messages = {
					{
						name = "BashExec",
						args = {
							[===[
              focus_path="${FM_FOCUS_PATH}"
              if [ "${focus_path}" ]; then
                nvim "${focus_path}"
              fi

              echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
              ]===],
						},
					},
				},
			},
			t = {
				help = "open new window in tmux",
				messages = {
					{
						name = "BashExec",
						args = {
							[===[
              focus_path="${FM_FOCUS_PATH}"
              if [ "${focus_path}" ]; then
                if [ -f "${focus_path}" ]; then
                  tmux new-window -c "$(dirname ${focus_path})"
                else
                  tmux new-window -c "${focus_path}"
                fi
              fi

              echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
              ]===],
						},
					},
				},
			},
			o = {
				help = "open",
				messages = {
					{
						name = "BashExec",
						args = {
							[===[
              focus_path="${FM_FOCUS_PATH}"
              if [ "${focus_path}" ]; then
                case $(file --mime-type ${focus_path} -b) in
                    text/*) nvim ${focus_path};;
                    application/json) nvim ${focus_path};;
                    *) open "${focus_path}"
                esac
              fi

              echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
              ]===],
						},
					},
				},
			},
			["ctrl+c"] = {
				help = "quit",
				messages = {
					{
						name = "Quit",
					},
				},
			},
		},
	},
}

fm.modes.builtins.default.keyBindings.onKeys["G"] = {
	help = "focus last",
	messages = {
		{
			name = "FocusLast",
		},
	},
}

fm.modes.builtins.default.keyBindings.onKeys["g"] = {
	help = "go to",
	messages = {
		{
			name = "SwitchMode",
			args = {
				"go-to",
			},
		},
	},
}

fm.modes.builtins.default.keyBindings.onKeys["o"] = {
	help = "open",
	messages = {
		{
			name = "SwitchMode",
			args = {
				"open",
			},
		},
	},
}

fm.modes.builtins.default.keyBindings.onKeys["y"] = {
	help = "yarn",
	messages = {
		{
			name = "SwitchMode",
			args = {
				"yarn",
			},
		},
	},
}

fm.modes.builtins.default.keyBindings.onKeys["/"] = {
	help = "search",
	messages = {
		{
			name = "BashExec",
			args = {
				[===[
        file_path=$(ls -a | fzf --no-sort)
        if [ "${file_path}" ]; then
          echo FocusPath "'"$PWD/${file_path}"'" >> "${FM_PIPE_MSG_IN:?}"
        fi
        ]===],
			},
		},
	},
}

fm.modes.builtins.default.keyBindings.onKeys["m"] = {
	help = "mark save",
	messages = {
		{
			name = "SwitchMode",
			args = {
				"mark-save",
			},
		},
	},
}

fm.modes.builtins.default.keyBindings.onKeys["`"] = {
	help = "mark load",
	messages = {
		{
			name = "SwitchMode",
			args = {
				"mark-load",
			},
		},
	},
}
