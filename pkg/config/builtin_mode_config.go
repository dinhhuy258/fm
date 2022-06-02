package config

var defaultModeConfig = ModeConfig{
	Name: "default",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Messages: []*MessageConfig{
					{
						Name: "Quit",
					},
				},
			},
			"j": {
				Help: "down",
				Messages: []*MessageConfig{
					{
						Name: "FocusNext",
					},
				},
			},
			"k": {
				Help: "up",
				Messages: []*MessageConfig{
					{
						Name: "FocusPrevious",
					},
				},
			},
			"l": {
				Help: "enter",
				Messages: []*MessageConfig{
					{
						Name: "Enter",
					},
				},
			},
			"h": {
				Help: "back",
				Messages: []*MessageConfig{
					{
						Name: "Back",
					},
				},
			},
			"m": {
				Help: "mark save",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"mark-save"},
					},
				},
			},
			"`": {
				Help: "mark load",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"mark-load"},
					},
				},
			},
			"n": {
				Help: "new file",
				Messages: []*MessageConfig{
					{
						Name: "SetInputBuffer",
						Args: []string{""},
					},
					{
						Name: "SwitchMode",
						Args: []string{"new-file"},
					},
				},
			},
			"s": {
				Help: "sort",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"sort"},
					},
				},
			},
			"r": {
				Help: "rename",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"rename"},
					},
					{
						Name: "BashExecSilently",
						Args: []string{`
              echo SetInputBuffer "'"$(basename "${FM_FOCUS_PATH}")"'" >> "${FM_PIPE_MSG_IN:?}"
						`},
					},
				},
			},
			"d": {
				Help: "delete",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"delete"},
					},
				},
			},
			"p": {
				Help: "copy",
				Messages: []*MessageConfig{
					{
						Name: "BashExec",
						Args: []string{`
							(while IFS= read -r line; do
								if cp -r -- "${line:?}" ./; then
									echo "${line:?}" copied to $PWD
								else
									echo Failed to copy "${line:?}" to $PWD
								fi
              done < "${FM_PIPE_SELECTION:?}")

							read -p "[Press enter to continue]"

							echo ClearSelection >> "${FM_PIPE_MSG_IN:?}"
							echo Refresh >> "${FM_PIPE_MSG_IN:?}"
						`},
					},
				},
			},
			"x": {
				Help: "move",
				Messages: []*MessageConfig{
					{
						Name: "BashExec",
						Args: []string{`
							(while IFS= read -r line; do
								if mv -- "${line:?}" ./; then
									echo "${line:?}" moved to $PWD
								else
									echo Failed to move "${line:?}" to $PWD
								fi
              done < "${FM_PIPE_SELECTION:?}")

							read -p "[Press enter to continue]"

							echo ClearSelection >> "${FM_PIPE_MSG_IN:?}"
							echo Refresh >> "${FM_PIPE_MSG_IN:?}"
						`},
					},
				},
			},
			"ctrl+r": {
				Help: "refresh",
				Messages: []*MessageConfig{
					{
						Name: "Refresh",
					},
				},
			},
			"space": {
				Help: "toggle selection",
				Messages: []*MessageConfig{
					{
						Name: "ToggleSelection",
					},
				},
			},
			"ctrl+space": {
				Help: "clear selection",
				Messages: []*MessageConfig{
					{
						Name: "ClearSelection",
					},
				},
			},
			".": {
				Help: "toggle hidden",
				Messages: []*MessageConfig{
					{
						Name: "ToggleHidden",
					},
				},
			},
		},
	},
}

var markSaveModeConfig = ModeConfig{
	Name: "mark-save",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Messages: []*MessageConfig{
					{
						Name: "Quit",
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "PopMode",
					},
				},
			},
		},
		Default: &ActionConfig{
			Messages: []*MessageConfig{
				{
					Name: "MarkSave",
				},
				{
					Name: "PopMode",
				},
			},
		},
	},
}

var markLoadModeConfig = ModeConfig{
	Name: "mark-load",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Messages: []*MessageConfig{
					{
						Name: "Quit",
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "PopMode",
					},
				},
			},
		},
		Default: &ActionConfig{
			Messages: []*MessageConfig{
				{
					Name: "MarkLoad",
				},
				{
					Name: "PopMode",
				},
			},
		},
	},
}

var newFileModeConfig = ModeConfig{
	Name: "new-file",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Messages: []*MessageConfig{
					{
						Name: "Quit",
					},
				},
			},
			"enter": {
				Help: "new file",
				Messages: []*MessageConfig{
					{
						Name: "BashExecSilently",
						Args: []string{`
							focus_path="${FM_FOCUS_PATH:?}"
							forcus_dir=$(dirname "$focus_path")
              name="${FM_INPUT_BUFFER}"

							if [[ "${name}" && ${name} == */ ]] ; then
							  name=${name%?}
								if [ -z "${name}" ]; then
									echo PopMode >> "${FM_PIPE_MSG_IN:?}"
								else
									mkdir -p -- "${name:?}" \
									&& echo Refresh >> "${FM_PIPE_MSG_IN:?}" \
									&& echo LogSuccess "'"${name} created"'" >> "${FM_PIPE_MSG_IN:?}" \
									&& echo FocusPath "'"${forcus_dir}/${name}"'" >> "${FM_PIPE_MSG_IN:?}" \
									&& echo PopMode >> "${FM_PIPE_MSG_IN:?}"
								fi
							elif [[ "${name}" ]] ; then
								mkdir -p -- "$(dirname ${name})" \
								&& touch -- "${name}" \
								&& echo Refresh >> "${FM_PIPE_MSG_IN:?}" \
								&& echo LogSuccess "'"${name} created"'" >> "${FM_PIPE_MSG_IN:?}" \
								&& echo FocusPath "'"${forcus_dir}/${name}"'" >> "${FM_PIPE_MSG_IN:?}" \
								&& echo PopMode >> "${FM_PIPE_MSG_IN:?}"
							else
								echo PopMode >> "${FM_PIPE_MSG_IN:?}"
              fi
						`},
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "PopMode",
					},
				},
			},
		},
		Default: &ActionConfig{
			Messages: []*MessageConfig{
				{
					Name: "UpdateInputBufferFromKey",
				},
			},
		},
	},
}

var renameModeConfig = ModeConfig{
	Name: "rename",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Messages: []*MessageConfig{
					{
						Name: "Quit",
					},
				},
			},
			"enter": {
				Help: "rename",
				Messages: []*MessageConfig{
					{
						Name: "BashExecSilently",
						Args: []string{`
							focus_path="${FM_FOCUS_PATH:?}"
              new_name="${FM_INPUT_BUFFER}"

							if [ -z "${new_name}" ]; then
								echo PopMode >> "${FM_PIPE_MSG_IN:?}"
							elif [ -e "${new_name:?}" ]; then
                echo LogError "'"${new_name} already exists"'" >> "${FM_PIPE_MSG_IN:?}"
								echo PopMode >> "${FM_PIPE_MSG_IN:?}"
              else
                mv -- "${focus_path:?}" "${new_name:?}" \
                  && echo Refresh >> "${FM_PIPE_MSG_IN:?}" \
                  && echo FocusPath "'"$(dirname "${focus_path}")/${new_name}"'" >> "${FM_PIPE_MSG_IN:?}" \
                  && echo LogSuccess "'"$(basename "${focus_path}") renamed to ${new_name}"'" >> "${FM_PIPE_MSG_IN:?}" \
									&& echo PopMode >> "${FM_PIPE_MSG_IN:?}"
              fi
						`},
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "PopMode",
					},
				},
			},
		},
		Default: &ActionConfig{
			Messages: []*MessageConfig{
				{
					Name: "UpdateInputBufferFromKey",
				},
			},
		},
	},
}

var deleteModeConfig = ModeConfig{
	Name: "delete",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Messages: []*MessageConfig{
					{
						Name: "Quit",
					},
				},
			},
			"d": {
				Help: "delete current",
				Messages: []*MessageConfig{
					{
						Name: "SetInputBuffer",
						Args: []string{"Do you want to delete this file? (y/n) "},
					},
					{
						Name: "SwitchMode",
						Args: []string{"delete-current"},
					},
				},
			},
			"s": {
				Help: "delete selections",
				Messages: []*MessageConfig{
					{
						Name: "SetInputBuffer",
						Args: []string{"Do you want to delete selected files? (y/n) "},
					},
					{
						Name: "SwitchMode",
						Args: []string{"delete-selections"},
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "PopMode",
					},
				},
			},
		},
	},
}

var deleteCurrentModeConfig = ModeConfig{
	Name: "delete-current",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Messages: []*MessageConfig{
					{
						Name: "Quit",
					},
				},
			},
			"y": {
				Help: "delete",
				Messages: []*MessageConfig{
					{
						Name: "BashExec",
						Args: []string{`
							focus_path="${FM_FOCUS_PATH:?}"

							if rm -rf -- "${focus_path}"; then
								echo "${focus_path}" deleted
							else
								echo Failed to delete "${focus_path}"
							fi

							echo Refresh >> "${FM_PIPE_MSG_IN:?}"
							echo PopMode >> "${FM_PIPE_MSG_IN:?}"
							echo PopMode >> "${FM_PIPE_MSG_IN:?}"

							read -p "[Press enter to continue]"
						`},
					},
				},
			},
		},
		Default: &ActionConfig{
			Help: "cancel",
			Messages: []*MessageConfig{
				{
					Name: "PopMode",
				},
				{
					Name: "PopMode",
				},
			},
		},
	},
}

var deleteSelectionsModeConfig = ModeConfig{
	Name: "delete-selections",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Messages: []*MessageConfig{
					{
						Name: "Quit",
					},
				},
			},
			"y": {
				Help: "delete selections",
				Messages: []*MessageConfig{
					{
						Name: "BashExec",
						Args: []string{`
							(while IFS= read -r line; do
								if rm -rf -- "${line:?}"; then
									echo "${line:?}" deleted
								else
									echo Failed to delete "${line:?}"
								fi
              done < "${FM_PIPE_SELECTION:?}")

							read -p "[Press enter to continue]"

							echo ClearSelection >> "${FM_PIPE_MSG_IN:?}"
							echo Refresh >> "${FM_PIPE_MSG_IN:?}"
							echo PopMode >> "${FM_PIPE_MSG_IN:?}"
							echo PopMode >> "${FM_PIPE_MSG_IN:?}"
						`},
					},
				},
			},
		},
		Default: &ActionConfig{
			Help: "cancel",
			Messages: []*MessageConfig{
				{
					Name: "PopMode",
				},
				{
					Name: "PopMode",
				},
			},
		},
	},
}

var sortModeConfig = ModeConfig{
	Name: "sort",
	KeyBindings: KeyBindingsConfig{
		OnKeys: map[string]*ActionConfig{
			"ctrl+c": {
				Help: "quit",
				Messages: []*MessageConfig{
					{
						Name: "Quit",
					},
					{
						Name: "PopMode",
					},
				},
			},
			"d": {
				Help: "dir first",
				Messages: []*MessageConfig{
					{
						Name: "SortByDirFirst",
					},
					{
						Name: "PopMode",
					},
				},
			},
			"c": {
				Help: "date modified",
				Messages: []*MessageConfig{
					{
						Name: "SortByDateModified",
					},
					{
						Name: "PopMode",
					},
				},
			},
			"n": {
				Help: "name",
				Messages: []*MessageConfig{
					{
						Name: "SortByName",
					},
					{
						Name: "PopMode",
					},
				},
			},
			"s": {
				Help: "size",
				Messages: []*MessageConfig{
					{
						Name: "SortBySize",
					},
					{
						Name: "PopMode",
					},
				},
			},
			"e": {
				Help: "extension",
				Messages: []*MessageConfig{
					{
						Name: "SortByExtension",
					},
					{
						Name: "PopMode",
					},
				},
			},
			"r": {
				Help: "reverse",
				Messages: []*MessageConfig{
					{
						Name: "ReverseSort",
					},
					{
						Name: "PopMode",
					},
				},
			},
		},
		Default: &ActionConfig{
			Messages: []*MessageConfig{
				{
					Name: "PopMode",
				},
			},
		},
	},
}

var builtinModeConfigs = map[string]*ModeConfig{
	"default":           &defaultModeConfig,
	"mark-save":         &markSaveModeConfig,
	"mark-load":         &markLoadModeConfig,
	"new-file":          &newFileModeConfig,
	"rename":            &renameModeConfig,
	"delete":            &deleteModeConfig,
	"delete-current":    &deleteCurrentModeConfig,
	"delete-selections": &deleteSelectionsModeConfig,
	"sort":              &sortModeConfig,
}
