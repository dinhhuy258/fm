package config

// defaultModeConfig is the configuration for the default builtin mode.
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
			"J": {
				Help: "go to bottom",
				Messages: []*MessageConfig{
					{
						Name: "FocusLast",
					},
				},
			},
			"K": {
				Help: "go to top",
				Messages: []*MessageConfig{
					{
						Name: "FocusFirst",
					},
				},
			},
			"n": {
				Help: "new file",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"new-file"},
					},
					{
						Name: "SetInputBuffer",
						Args: []string{""},
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
						Name: "BashExec",
						Args: []string{`
							if [ -s "${FM_PIPE_SELECTION:?}" ]; then
								# Handle selection deletion
								read -n 1 -p "Do you want to delete selected files? (y/n) " confirmation
								echo  # move to a new line after keypress
								if [ "$confirmation" = "y" ]; then
									deleted_count=0
									(while IFS= read -r line; do
										if rm -rf -- "${line:?}"; then
											echo "${line:?}" deleted
											deleted_count=$((deleted_count + 1))
										else
											echo Failed to delete "${line:?}"
										fi
									done < "${FM_PIPE_SELECTION:?}")

									if [ "$deleted_count" -gt 0 ]; then
										success_msg="$deleted_count file(s) deleted successfully"
										echo LogSuccess "'"$success_msg"'" >> "${FM_PIPE_MSG_IN:?}"
									fi

									echo Refresh >> "${FM_PIPE_MSG_IN:?}"

									focus_path="${FM_FOCUS_PATH:?}"
									# Check if focus path is in FM_PIPE_SELECTION file
									if grep -q "${focus_path:?}" "${FM_PIPE_SELECTION:?}"; then
										focus_index="${FM_FOCUS_IDX:?}"
										echo FocusByIndex "'"$focus_index"'" >> "${FM_PIPE_MSG_IN:?}"
									else
										echo FocusPath "'"$focus_path"'" >> "${FM_PIPE_MSG_IN:?}"
									fi

									echo ClearSelection >> "${FM_PIPE_MSG_IN:?}"
								fi
							else
								# Handle current file deletion
								focus_path="${FM_FOCUS_PATH:?}"
								read -n 1 -p "Do you want to delete this file? (y/n) " confirmation
								echo  # move to a new line after keypress
								if [ "$confirmation" = "y" ]; then
									if rm -rf -- "${focus_path}"; then
										echo "${focus_path}" deleted
										echo LogSuccess "'"$(basename "${focus_path}") deleted successfully"'" >> "${FM_PIPE_MSG_IN:?}"
									else
										echo LogError "'"Failed to delete "${focus_path}""'" >> "${FM_PIPE_MSG_IN:?}"
									fi

									echo Refresh >> "${FM_PIPE_MSG_IN:?}"

									focus_index="${FM_FOCUS_IDX:?}"
									echo FocusByIndex "'"$focus_index"'" >> "${FM_PIPE_MSG_IN:?}"

									if grep -q "${focus_path:?}" "${FM_PIPE_SELECTION:?}"; then
										echo ToggleSelectionByPath "'"$focus_path"'" >> "${FM_PIPE_MSG_IN:?}"
									fi
								fi
							fi

							read -p "[Press enter to continue]"
						`},
					},
				},
			},
			"p": {
				Help: "copy",
				Messages: []*MessageConfig{
					{
						Name: "BashExec",
						Args: []string{`
							if [ -s "${FM_PIPE_SELECTION:?}" ]; then
								# Handle selection copy
								read -n 1 -p "Do you want to copy selected files here? (y/n) " confirmation
								echo  # move to a new line after keypress
								if [ "$confirmation" = "y" ]; then
									copied_count=0
									while IFS= read -r line; do
										if cp -r -- "${line:?}" ./; then
											copied_count=$((copied_count + 1))
										fi
									done < "${FM_PIPE_SELECTION:?}"

									echo Refresh >> "${FM_PIPE_MSG_IN:?}"
									echo ClearSelection >> "${FM_PIPE_MSG_IN:?}"

									if [ "$copied_count" -gt 0 ]; then
										success_msg="$copied_count file(s) copied successfully"
										echo LogSuccess "'"$success_msg"'" >> "${FM_PIPE_MSG_IN:?}"
									fi
								fi
							else
								echo LogWarning "'"Select nothing"'" >> "${FM_PIPE_MSG_IN:?}"
							fi
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
							if [ -s "${FM_PIPE_SELECTION:?}" ]; then
								# Handle selection move
								read -n 1 -p "Do you want to move selected files here? (y/n) " confirmation
								echo  # move to a new line after keypress
								if [ "$confirmation" = "y" ]; then
									moved_count=0
									while IFS= read -r line; do
										if mv -- "${line:?}" ./; then
											moved_count=$((moved_count + 1))
										fi
									done < "${FM_PIPE_SELECTION:?}"

									echo Refresh >> "${FM_PIPE_MSG_IN:?}"
									echo ClearSelection >> "${FM_PIPE_MSG_IN:?}"

									if [ "$moved_count" -gt 0 ]; then
										success_msg="$moved_count file(s) moved successfully"
										echo LogSuccess "'"$success_msg"'" >> "${FM_PIPE_MSG_IN:?}"
									fi
								fi
							else
								echo LogWarning "'"Select nothing"'" >> "${FM_PIPE_MSG_IN:?}"
							fi
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
			" ": {
				Help: "toggle selection",
				Messages: []*MessageConfig{
					{
						Name: "ToggleSelection",
					},
				},
			},
			"ctrl+@": {
				Help: "clear selection",
				Messages: []*MessageConfig{
					{
						Name: "ClearSelection",
					},
				},
			},
			"ctrl+a": {
				Help: "select all",
				Messages: []*MessageConfig{
					{
						Name: "SelectAll",
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
			":": {
				Help: "command",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"command"},
					},
					{
						Name: "SetInputBuffer",
						Args: []string{""},
					},
				},
			},
		},
		OnNumber: &ActionConfig{
			Messages: []*MessageConfig{
				{
					Name: "SwitchMode",
					Args: []string{"go-to-index"},
				},
				{
					Name: "SetInputBuffer",
					Args: []string{""},
				},
				{
					Name: "UpdateInputBufferFromKey",
				},
			},
		},
	},
}

// newFileModeConfig is the configuration for the new file builtin mode.
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
									echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
								else
									mkdir -p -- "${name:?}" \
									&& echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}" \
									&& echo Refresh >> "${FM_PIPE_MSG_IN:?}" \
									&& echo LogSuccess "'"${name} created"'" >> "${FM_PIPE_MSG_IN:?}" \
									&& echo FocusPath "'"${forcus_dir}/${name}"'" >> "${FM_PIPE_MSG_IN:?}"
								fi
							elif [[ "${name}" ]] ; then
								mkdir -p -- "$(dirname ${name})" \
								&& touch -- "${name}" \
								&& echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}" \
								&& echo Refresh >> "${FM_PIPE_MSG_IN:?}" \
								&& echo LogSuccess "'"${name} created"'" >> "${FM_PIPE_MSG_IN:?}" \
								&& echo FocusPath "'"${forcus_dir}/${name}"'" >> "${FM_PIPE_MSG_IN:?}"
							else
								echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
              fi
						`},
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"default"},
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

// renameModeConfig is the configuration for the rename builtin mode.
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
								echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
							elif [ -e "${new_name:?}" ]; then
								echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
                echo LogError "'"${new_name} already exists"'" >> "${FM_PIPE_MSG_IN:?}"
              else
                mv -- "${focus_path:?}" "${new_name:?}" \
                  && echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}" \
                  && echo Refresh >> "${FM_PIPE_MSG_IN:?}" \
                  && echo FocusPath "'"$(dirname "${focus_path}")/${new_name}"'" >> "${FM_PIPE_MSG_IN:?}" \
                  && echo LogSuccess "'"$(basename "${focus_path}") renamed to ${new_name}"'" >> "${FM_PIPE_MSG_IN:?}"
              fi
						`},
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"default"},
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

// deleteModeConfig is the configuration for the delete builtin mode.
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
						Name: "SwitchMode",
						Args: []string{"delete-current"},
					},
					{
						Name: "SetInputBuffer",
						Args: []string{"Do you want to delete this file? (y/n) "},
					},
				},
			},
			"s": {
				Help: "delete selections",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"delete-selections"},
					},
					{
						Name: "SetInputBuffer",
						Args: []string{"Do you want to delete selected files? (y/n) "},
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"default"},
					},
				},
			},
		},
	},
}

// sortModeConfig is the configuration for the sort builtin mode.
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
				},
			},
			"d": {
				Help: "dir first",
				Messages: []*MessageConfig{
					{
						Name: "SortByDirFirst",
					},
					{
						Name: "SwitchMode",
						Args: []string{"default"},
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
						Name: "SwitchMode",
						Args: []string{"default"},
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
						Name: "SwitchMode",
						Args: []string{"default"},
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
						Name: "SwitchMode",
						Args: []string{"default"},
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
						Name: "SwitchMode",
						Args: []string{"default"},
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
						Name: "SwitchMode",
						Args: []string{"default"},
					},
				},
			},
		},
		Default: &ActionConfig{
			Messages: []*MessageConfig{
				{
					Name: "SwitchMode",
					Args: []string{"default"},
				},
			},
		},
	},
}

// commandModeConfig is the configuration for the command builtin mode.
var commandModeConfig = ModeConfig{
	Name: "command",
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
				Help: "execute",
				Messages: []*MessageConfig{
					{
						Name: "BashExec",
						Args: []string{`
              command="${FM_INPUT_BUFFER}"
							eval "$command"

							read -p "[Press enter to continue]"

							echo Refresh >> "${FM_PIPE_MSG_IN:?}"
							echo SwitchMode "'"default"'" >> "${FM_PIPE_MSG_IN:?}"
						`},
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"default"},
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

// goToIndexModeConfig is the configuration for the go to index builtin mode.
var goToIndexModeConfig = ModeConfig{
	Name: "go-to-index",
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
				Help: "go to index",
				Messages: []*MessageConfig{
					{
						Name: "BashExecSilently",
						Args: []string{`
							focus_index="${FM_INPUT_BUFFER}"
							echo FocusByIndex "'"$((focus_index-1))"'" >> "${FM_PIPE_MSG_IN:?}"
						`},
					},
					{
						Name: "SwitchMode",
						Args: []string{"default"},
					},
				},
			},
			"esc": {
				Help: "cancel",
				Messages: []*MessageConfig{
					{
						Name: "SwitchMode",
						Args: []string{"default"},
					},
				},
			},
			"backspace": {
				Messages: []*MessageConfig{
					{
						Name: "UpdateInputBufferFromKey",
					},
				},
			},
		},
		OnNumber: &ActionConfig{
			Messages: []*MessageConfig{
				{
					Name: "UpdateInputBufferFromKey",
				},
			},
		},
		Default: &ActionConfig{
			Messages: []*MessageConfig{
				{
					Name: "Null",
				},
			},
		},
	},
}

// builtinModeConfigs is a map of mode names to their configs.
var builtinModeConfigs = map[string]*ModeConfig{
	"default":     &defaultModeConfig,
	"new-file":    &newFileModeConfig,
	"rename":      &renameModeConfig,
	"sort":        &sortModeConfig,
	"command":     &commandModeConfig,
	"go-to-index": &goToIndexModeConfig,
}
