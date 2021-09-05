package fs

import (
	"os"
)

type FileManager struct {
	directory *Directory
}

func NewFileManager() (*FileManager, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fileManager := &FileManager{}
	fileManager.loadDirectory(wd)

	return fileManager, nil
}

func (fm *FileManager) loadDirectory(path string) {
	fm.directory = NewDirectory(path)
}
