package fs

import (
	"log"
	"os"
)

type FileManager struct {
	Dir             *Directory
	DirLoadedChan   chan struct{}
	DeleteCountChan chan int
	DeleteErrChan   chan error
}

func NewFileManager() (*FileManager, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fileManager := &FileManager{
		DirLoadedChan:   make(chan struct{}, 1),
		DeleteCountChan: make(chan int, 1024),
		DeleteErrChan:   make(chan error),
	}
	fileManager.LoadDirectory(wd)

	return fileManager, nil
}

func (fm *FileManager) LoadDirectory(path string) {
	fm.Dir = &Directory{
		Path: path,
	}

	fm.Reload()
}

func (fm *FileManager) Reload() {
	go func() {
		err := fm.Dir.ReadDir()
		if err != nil {
			log.Printf("failed to read directory %v", err)
		}
		fm.DirLoadedChan <- struct{}{}
	}()
}

func (fm *FileManager) Delete(paths []string) {
	go func() {
		for _, path := range paths {
			if err := os.RemoveAll(path); err != nil {
				fm.DeleteErrChan <- err
			}

			fm.DeleteCountChan <- 1
		}
	}()
}
