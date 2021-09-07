package fs

import (
	"log"
	"os"
)

type FileManager struct {
	Dir           *Directory
	DirLoadedChan chan struct{}
}

func NewFileManager() (*FileManager, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fileManager := &FileManager{
		DirLoadedChan: make(chan struct{}, 1),
	}
	fileManager.LoadDirectory(wd)

	return fileManager, nil
}

func (fm *FileManager) LoadDirectory(path string) {
	fm.Dir = &Directory{
		Path: path,
	}

	go func() {
		err := fm.Dir.ReadDir()
		if err != nil {
			log.Printf("failed to read directory %v", err)
		}
		fm.DirLoadedChan <- struct{}{}
	}()
}
