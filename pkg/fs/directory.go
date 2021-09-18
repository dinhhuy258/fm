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
	Nodes        []*Node
	VisibleNodes []*Node
	Path         string
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
	if err := f.Close(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	cfg := config.AppConfig
	nodes := make([]*Node, 0, len(names))
	visibleNodes := make([]*Node, 0, len(names))

	for _, relativePath := range names {
		absolutePath := filepath.Join(dir.Path, relativePath)
		lstat, err := os.Lstat(absolutePath)

		if os.IsNotExist(err) {
			continue
		}

		if err != nil {
			return err
		}

		node := &Node{
			RelativePath: relativePath,
			AbsolutePath: absolutePath,
			IsDir:        lstat.IsDir(),
			Size:         lstat.Size(),
			Extension:    filepath.Ext(absolutePath),
		}

		if cfg.ShowHidden || !isHidden(relativePath) {
			visibleNodes = append(visibleNodes, node)
		}

		nodes = append(nodes, node)
	}

	dir.VisibleNodes = visibleNodes
	dir.Nodes = nodes

	return nil
}

func (dir *Directory) Reload() {
	visibleNodes := make([]*Node, 0, len(dir.Nodes))
	cfg := config.AppConfig

	for _, node := range dir.Nodes {
		if cfg.ShowHidden || !isHidden(node.RelativePath) {
			visibleNodes = append(visibleNodes, node)
		}
	}

	dir.VisibleNodes = visibleNodes
}
