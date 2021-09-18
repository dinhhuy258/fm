package fs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileManager struct {
	Dir           *Directory
	DirLoadedChan chan struct{}
	OpCountChan   chan int
	OpErrChan     chan error
}

func NewFileManager() (*FileManager, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fileManager := &FileManager{
		DirLoadedChan: make(chan struct{}, 1),
		OpCountChan:   make(chan int, 1024),
		OpErrChan:     make(chan error),
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
				fm.OpErrChan <- err
			}

			fm.OpCountChan <- 1
		}
	}()
}

func (fm *FileManager) Copy(
	paths []string,
	destDir string,
) (countChan chan int64, errChan chan error) {
	countChan = make(chan int64, len(paths))
	errChan = make(chan error, len(paths))

	go func() {
		for _, path := range paths {
			dst := filepath.Join(destDir, filepath.Base(path))
			_, err := os.Lstat(dst)

			if !os.IsNotExist(err) {
				var newPath string

				for i := 1; !os.IsNotExist(err); i++ {
					newPath = fmt.Sprintf("%s.~%d~", dst, i)
					_, err = os.Lstat(newPath)
				}

				dst = newPath
			}

			walkFunc := func(path string, info os.FileInfo, err error) error {
				if err != nil {
					errChan <- err

					return nil
				}

				if err := copyPath(path, dst, info); err != nil {
					errChan <- err
				}

				return nil
			}

			if err := filepath.Walk(path, walkFunc); err != nil {
				errChan <- fmt.Errorf("walk: %w", err)
			}

			countChan <- 1
		}
	}()

	return countChan, errChan
}
