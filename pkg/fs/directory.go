package fs

import (
	"os"
	"path/filepath"

	"github.com/dinhhuy258/fm/pkg/config"
)

type Node struct {
	RelativePath string
	AbsolutePath string
	Extension    string
	IsDir        bool
	Size         int64
}

type Directory struct {
	Nodes []*Node
	Path  string
}

func (dir *Directory) Parent() string {
	return filepath.Dir(dir.Path)
}

func (dir *Directory) ReadDir() error {
	f, err := os.Open(dir.Path)
	if err != nil {
		return err
	}

	names, err := f.Readdirnames(-1)
	f.Close()

	if err != nil {
		return err
	}

	config := config.AppConfig
	nodes := make([]*Node, 0, len(names))

	for _, relativePath := range names {
		absolutePath := filepath.Join(dir.Path, relativePath)
		lstat, err := os.Lstat(absolutePath)

		if os.IsNotExist(err) || (!config.ShowHidden && isHidden(relativePath)) {
			continue
		}

		if err != nil {
			return err
		}

		nodes = append(nodes,
			&Node{
				RelativePath: relativePath,
				AbsolutePath: absolutePath,
				IsDir:        lstat.IsDir(),
				Size:         lstat.Size(),
				Extension:    filepath.Ext(absolutePath),
			},
		)
	}

	dir.Nodes = nodes

	return nil
}
