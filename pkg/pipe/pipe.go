package pipe

import (
	"os"
	"path/filepath"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/hpcloud/tail"
)

// Pipe is a pipe to communicate with fm
type Pipe struct {
	sessionPath      string
	messageInPath    string
	selectionPath    string
	messageInWatcher *tail.Tail
	watcherStop      chan bool
}

// NewPipe creates a new pipe
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
		messageInPath:    messageInPath,
		selectionPath:    selectionPath,
		messageInWatcher: messageInWatcher,
		watcherStop:      make(chan bool),
	}, nil
}

// GetSessionPath returns session the path of the application
func (p *Pipe) GetSessionPath() string {
	return p.sessionPath
}

// GetMessageInPath returns the path to the message in file
func (p *Pipe) GetMessageInPath() string {
	return p.messageInPath
}

// GetSelectionPath returns the path to the selection file
func (p *Pipe) GetSelectionPath() string {
	return p.selectionPath
}

// StartWatcher starts the watcher for message in file
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

// StopWatcher stops the watcher for message in file
func (p *Pipe) StopWatcher() {
	p.watcherStop <- true
	p.messageInWatcher.Cleanup()
	_ = p.messageInWatcher.Stop()
}
