package fs

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/dinhhuy258/fm/pkg/config"
)

type FileExplorer struct {
	path    string
	entries []IEntry
}

var (
	fileExplorer             *FileExplorer
	fileExplorerCreationOnce sync.Once
)

func GetFileExplorer() *FileExplorer {
	fileExplorerCreationOnce.Do(func() {
		fileExplorer = &FileExplorer{}
	})

	return fileExplorer
}

func (fe *FileExplorer) GetPath() string {
	return fe.path
}

func (fe *FileExplorer) Dir() string {
	return Dir(fe.path)
}

func (fe *FileExplorer) GetEntriesSize() int {
	return len(fe.entries)
}

func (fe *FileExplorer) GetEntries() []IEntry {
	return fe.entries
}

func (fe *FileExplorer) GetEntry(idx int) IEntry {
	if idx < 0 || idx >= len(fe.entries) {
		log.Fatalf("invalid idx: %d", idx)
	}

	return fe.entries[idx]
}

func (fe *FileExplorer) LoadEntries(path string, onEntriesLoaded func()) {
	//TODO: Check if path is directory or not
	fe.path = path

	go func() {
		if err := fe.loadEntries(); err != nil {
			log.Fatalf("failed to load entries for the current path: %s", fe.path)
		}

		onEntriesLoaded()
	}()
}

func (fe *FileExplorer) loadEntries() error {
	f, err := os.Open(fe.path)
	if err != nil {
		return err
	}

	names, err := f.Readdirnames(-1)
	if err := f.Close(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	cfg := config.AppConfig
	fe.entries = make([]IEntry, 0, len(names))

	for _, name := range names {
		absolutePath := filepath.Join(fe.path, name)
		lstat, err := os.Lstat(absolutePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return err
		}

		if !cfg.ShowHidden && isHidden(name) {
			continue
		}

		isDir := lstat.IsDir()
		size := lstat.Size()

		if isDir {
			fe.entries = append(fe.entries, &Directory2{
				&Entry{
					name: name,
					path: absolutePath,
					size: size,
				},
			})
		} else {
			fe.entries = append(fe.entries, &File{
				&Entry{
					name: name,
					path: absolutePath,
					size: size,
				},
			})
		}
	}

	return nil
}
