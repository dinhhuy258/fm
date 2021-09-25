package fs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileManager struct {
	Dir *Directory
}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (fm *FileManager) LoadDirectory(path string) chan struct{} {
	fm.Dir = &Directory{
		Path: path,
	}

	dirLoadedChan := make(chan struct{}, 1)

	go func() {
		err := fm.Dir.ReadDir()
		if err != nil {
			log.Fatalf("failed to read directory %s. Error: %v", fm.Dir.Path, err)
		}

		dirLoadedChan <- struct{}{}
	}()

	return dirLoadedChan
}

func (fm *FileManager) Delete(paths []string) (countChan chan int, errChan chan error) {
	countChan = make(chan int, len(paths))
	errChan = make(chan error)

	go func() {
		for _, path := range paths {
			if err := os.RemoveAll(path); err != nil {
				errChan <- err
			}

			countChan <- 1
		}
	}()

	return countChan, errChan
}

func (fm *FileManager) Copy(
	srcPaths []string,
	destDir string,
) (countChan chan int, errChan chan error) {
	countChan = make(chan int, len(srcPaths))
	errChan = make(chan error)

	go func() {
		for _, srcPath := range srcPaths {
			dst := filepath.Join(destDir, filepath.Base(srcPath))
			_, err := os.Lstat(dst)

			if !os.IsNotExist(err) {
				var newPath string

				for i := 1; !os.IsNotExist(err); i++ {
					newPath = fmt.Sprintf("%s.~%d~", dst, i)
					_, err = os.Lstat(newPath)
				}

				dst = newPath
			}

			src := srcPath // This will make scopelint happy
			walkFunc := func(path string, info os.FileInfo, err error) error {
				if err != nil {
					errChan <- err

					return nil
				}

				if err := copyPath(src, path, dst, info); err != nil {
					errChan <- err
				}

				return nil
			}

			if err := filepath.Walk(srcPath, walkFunc); err != nil {
				errChan <- err
			}

			countChan <- 1
		}
	}()

	return countChan, errChan
}

func (fm *FileManager) Move(
	srcPaths []string,
	destDir string,
) (countChan chan int, errChan chan error) {
	countChan = make(chan int, len(srcPaths))
	errChan = make(chan error)

	go func() {
		for _, src := range srcPaths {
			dst := filepath.Join(destDir, filepath.Base(src))
			if dst == src {
				countChan <- 1

				continue
			}

			_, err := os.Stat(dst)
			if !os.IsNotExist(err) {
				var newPath string

				for i := 1; !os.IsNotExist(err); i++ {
					newPath = fmt.Sprintf("%s.~%d~", dst, i)
					_, err = os.Lstat(newPath)
				}

				dst = newPath
			}

			if err := os.Rename(src, dst); err != nil {
				errChan <- err
			}

			countChan <- 1
		}
	}()

	return countChan, errChan
}
