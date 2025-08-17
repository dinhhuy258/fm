package fs

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/djherbis/times"

	"github.com/dinhhuy258/fm/pkg/types"
)

var errFileNotFound = errors.New("file not found")

// IEntry represents a file or directory.
type IEntry interface {
	GetName() string
	GetPath() string
	GetSize() int64
	GetExt() string
	GetPermissions() string
	IsDirectory() bool
	IsSymlink() bool
	GetChangeTime() time.Time
}

// Entry contains information about a file or directory.
type Entry struct {
	IEntry

	name        string
	path        string
	isSymlink   bool
	size        int64
	ext         string
	permissions string
	changeTime  time.Time
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

// GetChangeTime returns the change time of the entry.
func (e *Entry) GetChangeTime() time.Time {
	return e.changeTime
}

// IsSymlink returns true if the current file is symlink
func (e *Entry) IsSymlink() bool {
	return e.isSymlink
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
func LoadEntries(path string,
	showHidden bool,
	sortAlgorithm string,
	sortReverse bool,
	sortIgnoreCase bool,
	sortIgnoreDiacritics bool,
) ([]IEntry, error) {
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
		entry, err := loadEntry(path, name, showHidden)
		if err != nil {
			if errors.Is(err, errFileNotFound) {
				continue
			}

			return nil, err
		}

		entries = append(entries, entry)
	}

	getEntrySort(types.SortType(sortAlgorithm)).sort(entries, sortReverse, sortIgnoreCase, sortIgnoreDiacritics)

	return entries, nil
}

// loadEntry loads the entry of the given file path and file name.
func loadEntry(path, name string, showHidden bool) (IEntry, error) {
	fpath := filepath.Join(path, name)

	lstat, err := getFileInfo(fpath)
	if err != nil {
		return nil, err
	}

	if !showHidden && isHidden(name) {
		return nil, errFileNotFound
	}

	var ct time.Time

	ts := times.Get(lstat)
	if ts.HasChangeTime() {
		ct = ts.ChangeTime()
	} else {
		ct = lstat.ModTime()
	}

	isDir := lstat.IsDir()
	size := lstat.Size()
	permissions := lstat.Mode().String()[1:]

	isSymlink := (lstat.Mode() & os.ModeSymlink) != 0
	if isSymlink {
		linkTarget, err := evalSymlinks(fpath)
		if err != nil {
			return nil, err
		}

		linkTargetLstat, err := getFileInfo(linkTarget)
		if err != nil {
			return nil, err
		}

		isDir = linkTargetLstat.IsDir()
	}

	var ext string
	if ext = filepath.Ext(fpath); ext != "" {
		ext = ext[1:]
	}

	if isDir {
		return &Directory{
			&Entry{
				name:        name,
				path:        fpath,
				size:        size,
				permissions: permissions,
				ext:         ext,
				changeTime:  ct,
				isSymlink:   isSymlink,
			},
		}, nil
	}

	return &File{
		&Entry{
			name:        name,
			path:        fpath,
			size:        size,
			permissions: permissions,
			ext:         ext,
			changeTime:  ct,
			isSymlink:   isSymlink,
		},
	}, nil
}

// getFileInfo get file information from the given file path
func getFileInfo(fpath string) (os.FileInfo, error) {
	lstat, err := os.Lstat(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errFileNotFound
		}

		return nil, err
	}

	return lstat, nil
}

// evalSymlinks get link target of the given symlinks
func evalSymlinks(fpath string) (string, error) {
	linkTarget, err := filepath.EvalSymlinks(fpath)
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			return "", errFileNotFound
		}

		return "", err
	}

	return linkTarget, nil
}

// isHidden returns true if the given name is a hidden file.
func isHidden(filename string) bool {
	return filename[0:1] == "."
}