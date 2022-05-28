package pipe

import (
	"os"
	"path/filepath"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/hpcloud/tail"
)

type Pipe struct {
	sessionPath      string
	selectionPath    string
	messageInPath    string
	messageInWatcher *tail.Tail
	watcherStop      chan bool
}

func NewPipe() (*Pipe, error) {
	runtime := os.Getenv("XDG_RUNTIME_DIR")
	if runtime == "" {
		runtime = os.TempDir()
	}

	sessionPath := filepath.Join(runtime, "fm", "session")
	if err := os.MkdirAll(sessionPath, os.ModePerm); err != nil {
		return nil, err
	}

	messageInPath := filepath.Join(sessionPath, "msg_in")
	if err := fs.CreateFile(messageInPath, true); err != nil {
		return nil, err
	}

	selectionPath := filepath.Join(sessionPath, "sellection")
	if err := fs.CreateFile(selectionPath, true); err != nil {
		return nil, err
	}

	messageInWatcher, err := tail.TailFile(messageInPath, tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}

	return &Pipe{
		sessionPath:      sessionPath,
		selectionPath:    selectionPath,
		messageInPath:    messageInPath,
		messageInWatcher: messageInWatcher,
		watcherStop:      make(chan bool),
	}, nil
}

func (p *Pipe) GetMessageInPath() string {
	return p.messageInPath
}

func (p *Pipe) GetSelectionPath() string {
	return p.selectionPath
}

func (p *Pipe) StartWatcher(onMessageIn func(string)) {
	go func() {
		for {
			select {
			case <-p.watcherStop:
				return
			case line := <-p.messageInWatcher.Lines:
				onMessageIn(line.Text)
			}
		}
	}()
}

func (p *Pipe) StopWatcher() {
	p.watcherStop <- true
	p.messageInWatcher.Cleanup()
	_ = p.messageInWatcher.Stop()
}
