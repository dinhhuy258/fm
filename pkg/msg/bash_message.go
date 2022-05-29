package msg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/controller"
	"github.com/dinhhuy258/fm/pkg/gui/view"
)

func BashExec(app IApp, params ...string) {
	// This function should be called in a UI thread, otherwise it will not work
	app.OnUIThread(func() error {
		if err := app.Suspend(); err != nil {
			return err
		}
		defer func() {
			_ = app.Resume()
		}()

		// Clear the terminal screen first
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			return err
		}

		command := params[0]
		cmd = exec.Command("bash", "-c", command)
		cmd.Env = getEnv(app)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}

		cmd.Stdout = ioutil.Discard
		cmd.Stderr = ioutil.Discard
		cmd.Stdin = nil

		return nil
	})
}

func BashExecSilently(app IApp, params ...string) {
	logController, _ := app.GetController(controller.Log).(*controller.LogController)

	command := params[0]
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = getEnv(app)

	if err := cmd.Run(); err != nil {
		logController.SetLog(view.Error, "Failed to execute script")
		logController.UpdateView()
	}
}

func getEnv(app IApp) []string {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	inputController, _ := app.GetController(controller.Input).(*controller.InputController)
	selectionController, _ := app.GetController(controller.Sellection).(*controller.SelectionController)

	pipe := app.GetPipe()
	focusPath := ""

	currentEntry := explorerController.GetCurrentEntry()
	if currentEntry != nil {
		focusPath = currentEntry.GetPath()
	}

	// Write selected files to pipe
	fs.WriteToFile(pipe.GetSelectionPath(), selectionController.GetSelections(), true)

	env := os.Environ()
	env = append(env, fmt.Sprintf("FM_FOCUS_PATH=%s", focusPath))
	env = append(env, fmt.Sprintf("FM_INPUT_BUFFER=%s", inputController.GetInputBuffer()))
	env = append(env, fmt.Sprintf("FM_PIPE_MSG_IN=%s", pipe.GetMessageInPath()))
	env = append(env, fmt.Sprintf("FM_PIPE_SELECTION=%s", pipe.GetSelectionPath()))

	return env
}
