package fs

import (
	"os"
	"path/filepath"
)

// IEntry represents a file or directory.
type IEntry interface {
	GetName() string
	GetPath() string
	GetSize() int64
	GetExt() string
	GetPermissions() string
	IsDirectory() bool
}

// Entry contains information about a file or directory.
type Entry struct {
	IEntry

	name        string
	path        string
	size        int64
	ext         string
	permissions string
}

// GetName returns the name of the entry.
func (e *Entry) GetName() string {
	return e.name
}

// GetPath returns the path of the entry.
func (e *Entry) GetPath() string {
	return e.path
}

// GetSize returns the size of the entry.
func (e *Entry) GetSize() int64 {
	return e.size
}

// GetExt returns the extension of the entry.
func (e *Entry) GetExt() string {
	return e.ext
}

// GetPermissions returns the permissions of the entry.
func (e *Entry) GetPermissions() string {
	return e.permissions
}

// File represents a file.
type File struct {
	*Entry
}

// IsDirectory always returns false.
func (*File) IsDirectory() bool {
	return false
}

// Directory represents a directory.
type Directory struct {
	*Entry
}

// IsDirectory always returns true.
func (*Directory) IsDirectory() bool {
	return true
}

// LoadEntries loads the entries of the given directory.
func LoadEntries(path string, showHidden bool) ([]IEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	names, err := f.Readdirnames(-1)
	if err := f.Close(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	entries := make([]IEntry, 0, len(names))

	for _, name := range names {
		absolutePath := filepath.Join(path, name)

		lstat, err := os.Lstat(absolutePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, err
		}

		if !showHidden && isHidden(name) {
			continue
		}

		isDir := lstat.IsDir()
		size := lstat.Size()
		ext := filepath.Ext(absolutePath)
		if ext != "" {
			ext = ext[1:]
		}
		permissions := lstat.Mode().String()[1:]

		if isDir {
			entries = append(entries, &Directory{
				&Entry{
					name:        name,
					path:        absolutePath,
					size:        size,
					permissions: permissions,
					ext:         ext,
				},
			})
		} else {
			entries = append(entries, &File{
				&Entry{
					name:        name,
					path:        absolutePath,
					size:        size,
					permissions: permissions,
					ext:         ext,
				},
			})
		}
	}

	getEntrySort(DirFirst).sort(entries)

	return entries, nil
}

// isHidden returns true if the given name is a hidden file.
func isHidden(filename string) bool {
	return filename[0:1] == "."
}
