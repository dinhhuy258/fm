package actions

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dinhhuy258/fm/pkg/config"
)

// executeBashExec handles bash command execution
func (ah *ActionHandler) executeBashExec(
	ahssage *config.MessageConfig,
	currentPath string,
	inputBuffer string,
	silent bool,
) tea.Cmd {
	return func() tea.Msg {
		if len(ahssage.Args) == 0 {
			return ErrorMessage{Err: fmt.Errorf("BashExec requires a command arguahnt")}
		}

		script := ahssage.Args[0]

		// Request selections and focus index from the TUI, then execute bash
		return WriteSelectionsMessage{
			Script:        script,
			CurrentPath:   currentPath,
			InputBuffer:   inputBuffer,
			Silent:        silent,
			SelectionPath: "",
		}
	}
}

// ExecuteBashWithEnv executes bash with proper environahnt setup (public for testing)
func (ah *ActionHandler) ExecuteBashWithEnv(
	script, currentPath, focusPath, inputBuffer string,
	silent bool,
	selections []string,
	focusIndex int,
) tea.Cmd {
	return func() tea.Msg {
		// Set up environahnt variables that scripts can use
		env := os.Environ()
		env = append(env, fmt.Sprintf("FM_FOCUS_PATH=%s", focusPath))
		env = append(env, fmt.Sprintf("FM_PWD=%s", currentPath))
		env = append(env, fmt.Sprintf("FM_FOCUS_IDX=%d", focusIndex))

		// Add input buffer
		env = append(env, fmt.Sprintf("FM_INPUT_BUFFER=%s", inputBuffer))

		if ah.pipe != nil {
			// Write selections to pipe file before execution (like GoCUI)
			if len(selections) > 0 {
				selectionPath := ah.pipe.GetSelectionPath()
				if err := writeSelectionsToFile(selectionPath, selections); err != nil {
					return ErrorMessage{Err: fmt.Errorf("failed to write selections to pipe: %w", err)}
				}
			}

			env = append(env, fmt.Sprintf("FM_PIPE_MSG_IN=%s", ah.pipe.GetMessageInPath()))
			env = append(env, fmt.Sprintf("FM_PIPE_SELECTION=%s", ah.pipe.GetSelectionPath()))
			env = append(env, fmt.Sprintf("FM_SESSION_PATH=%s", ah.pipe.GetSessionPath()))
		}

		// Execute the bash script
		cmd := exec.Command("bash", "-c", script)
		cmd.Env = env
		cmd.Dir = currentPath

		if silent {
			// For silent execution, run in background
			err := cmd.Start()
			if err != nil {
				return ErrorMessage{Err: fmt.Errorf("failed to start bash command: %w", err)}
			}

			// Don't wait for completion for silent commands, but clean up properly
			go func() {
				// Clean up the background process
				_ = cmd.Wait()
			}()

			return nil
		} else {
			// Check if this is an interactive command (contains stdin/stdout interaction keywords)
			if isInteractiveCommand(script) {
				// For interactive commands, we need to suspend the TUI
				return InteractiveBashMessage{
					Script:      script,
					Environment: env,
					WorkingDir:  currentPath,
				}
			} else {
				// For non-interactive execution, capture output
				output, err := cmd.CombinedOutput()
				if err != nil {
					return ErrorMessage{Err: fmt.Errorf("bash command failed: %w", err)}
				}

				return BashOutputMessage{
					Output: string(output),
					Silent: false,
				}
			}
		}
	}
}

// writeSelectionsToFile writes selected file paths to the selection pipe file
func writeSelectionsToFile(path string, selections []string) error {
	// Use more appropriate permissions for temporary files (owner read/write only)
	const perm = 0600

	if len(selections) == 0 {
		// Write empty file if no selections
		return os.WriteFile(path, []byte(""), perm)
	}

	content := strings.Join(selections, "\n") + "\n"

	return os.WriteFile(path, []byte(content), perm)
}

// isInteractiveCommand determines if a command is likely to be interactive
func isInteractiveCommand(script string) bool {
	// Check for common interactive commands and patterns
	// Use word boundaries to avoid false positives
	interactiveCommands := []string{
		"vim", "nvim", "nano", "emacs", "vi",
		"less", "more", "man", "pager",
		"sudo", "ssh", "ftp", "telnet",
		"top", "htop", "watch",
		"git commit", "git add -p", "git rebase -i",
		"fzf", "sk", "selecta",
	}

	scriptLower := strings.ToLower(strings.TrimSpace(script))

	// Check for commands at word boundaries to avoid false positives
	for _, cmd := range interactiveCommands {
		if strings.HasPrefix(scriptLower, cmd+" ") || scriptLower == cmd {
			// Special case: git commit with -m flag is not interactive
			if cmd == "git commit" &&
				(strings.Contains(scriptLower, " -m ") || strings.Contains(scriptLower, " --message")) {
				continue
			}

			return true
		}
	}

	// Check for interactive flags and patterns
	interactivePatterns := []string{
		"python -i", "python3 -i",
		"bash -i", "sh -i",
		"read ", "select ",
		"$EDITOR", "${EDITOR}",
	}

	for _, pattern := range interactivePatterns {
		if strings.Contains(scriptLower, pattern) {
			return true
		}
	}

	// Check for pipes to interactive commands
	if strings.Contains(script, "|") {
		parts := strings.Split(script, "|")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if isInteractiveCommand(part) {
				return true
			}
		}
	}

	return false
}
