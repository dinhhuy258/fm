package message

import (
	// "io/ioutil"
	"io/ioutil"
	"os"
	"os/exec"
)

func BashExec(app IApp, params ...string) {
	app.Suspend()

	command := params[0]

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	cmd.Stdin = nil

	app.Resume()
}
