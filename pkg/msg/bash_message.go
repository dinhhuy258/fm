package msg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/dinhhuy258/fm/pkg/gui/controller"
)

func BashExec(app IApp, params ...string) {
	app.OnUIThread(func() error {
		if err := app.Suspend(); err != nil {
			return err
		}

		command := params[0]
		cmd := exec.Command("bash", "-c", command)
		cmd.Env = getEnv(app)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			_ = app.Resume()

			return err
		}

		cmd.Stdout = ioutil.Discard
		cmd.Stderr = ioutil.Discard
		cmd.Stdin = nil

		if err := app.Resume(); err != nil {
			return err
		}

		return nil
	})
}

func getEnv(app IApp) []string {
	explorerController, _ := app.GetController(controller.Explorer).(*controller.ExplorerController)
	currentEntry := explorerController.GetCurrentEntry()

	env := os.Environ()
	env = append(env, fmt.Sprintf("FM_FOCUS_PATH=%s", currentEntry.GetPath()))
	env = append(env, fmt.Sprintf("FM_PIPE_MSG_IN=%s", app.GetPipe().GetMsgInPath()))

	return env
}
