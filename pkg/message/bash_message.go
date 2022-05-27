package message

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func BashExec(app IApp, params ...string) {
	app.OnUIThread(func() error {
		if err := app.Suspend(); err != nil {
			return err
		}

		command := params[0]
		cmd := exec.Command("bash", "-c", command)
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