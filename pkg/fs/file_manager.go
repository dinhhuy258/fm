package fs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileManager struct {
	dir *Directory
}

func (fm *FileManager) IsLoaded() bool {
	return fm.dir != nil
}

func (fm *FileManager) Reload() {
	fm.dir.Reload()
}

func (fm *FileManager) GetNodeAtIdx(idx int) *Node {
	if idx < 0 || idx >= len(fm.dir.VisibleNodes) {
		log.Fatalf("invalid idx: %d", idx)
	}

	return fm.dir.VisibleNodes[idx]
}

func (fm *FileManager) GetCurrentPath() string {
	return fm.dir.Path
}

func (fm *FileManager) GetVisibleNodes() []*Node {
	return fm.dir.VisibleNodes
}

func (fm *FileManager) GetVisibleNodesSize() int {
	return len(fm.dir.VisibleNodes)
}

func (fm *FileManager) GetParentPath() string {
	return fm.dir.Parent()
}

var (
	fileManager             *FileManager
	fileManagerCreationOnce sync.Once
)

func GetFileManager() *FileManager {
	fileManagerCreationOnce.Do(func() {
		fileManager = &FileManager{}
	})

	return fileManager
}

func (fm *FileManager) LoadDirectory(path string) chan struct{} {
	fm.dir = &Directory{
		Path: path,
	}

	dirLoadedChan := make(chan struct{}, 1)

	go func() {
		err := fm.dir.ReadDir()
		if err != nil {
			log.Fatalf("failed to read directory %s. Error: %v", fm.dir.Path, err)
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
