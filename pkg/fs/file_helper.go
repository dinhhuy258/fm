package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Humanize(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%d B", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / 1000
	for _, s := range suffix {
		if curr < 10 {
			return fmt.Sprintf("%.1f %s", curr-0.0499, s)
		} else if curr < 1000 {
			return fmt.Sprintf("%d %s", int(curr), s)
		}

		curr /= 1000
	}

	return ""
}

func Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func CreateFile(name string, override bool) error {
	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	if override {
		flags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	}

	f, err := os.OpenFile(name, flags, 0o644)
	if err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return err
}

// CreateDirectory creates a new directory given a name.
func CreateDirectory(name string) error {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(name, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func Dir(path string) string {
	return filepath.Dir(path)
}

func Delete(paths []string, onSuccess func(), onError func(error), onComplete func(int, int)) {
	successCount := 0
	errorCount := 0

	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			errorCount++

			onError(err)
		} else {
			successCount++

			onSuccess()
		}
	}

	onComplete(successCount, errorCount)
}

func Copy(srcPaths []string, destDir string, onSuccess func(), onError func(error), onComplete func(int, int)) {
	successCount := 0
	errorCount := 0

	for _, srcPath := range srcPaths {
		srcPath := srcPath
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

		walkFunc := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			return copyPath(srcPath, path, dst, info)
		}

		if err := filepath.Walk(srcPath, walkFunc); err != nil {
			errorCount++

			onError(err)
		} else {
			successCount++

			onSuccess()
		}
	}

	onComplete(successCount, errorCount)
}

func Move(srcPaths []string, destDir string, onSuccess func(), onError func(error), onComplete func(int, int)) {
	successCount := 0
	errorCount := 0

	for _, src := range srcPaths {
		dst := filepath.Join(destDir, filepath.Base(src))
		if dst == src {
			successCount++

			onSuccess()

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
			errorCount++

			onError(err)
		} else {
			successCount++

			onSuccess()
		}
	}

	onComplete(successCount, errorCount)
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

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

		if isDir {
			entries = append(entries, &Directory{
				&Entry{
					name: name,
					path: absolutePath,
					size: size,
				},
			})
		} else {
			entries = append(entries, &File{
				&Entry{
					name: name,
					path: absolutePath,
					size: size,
				},
			})
		}
	}

	return entries, nil
}

func isHidden(filename string) bool {
	return filename[0:1] == "."
}

func copyFile(src, dst string, info os.FileInfo) error {
	buf := make([]byte, 4096)

	r, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func(r *os.File) {
		_ = r.Close()
	}(r)

	w, err := os.Create(dst)
	if err != nil {
		return err
	}

	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			_ = w.Close()
			_ = os.Remove(dst)

			return err
		}

		if n == 0 {
			break
		}

		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
	}

	if err := w.Close(); err != nil {
		_ = os.Remove(dst)

		return err
	}

	if err := os.Chmod(dst, info.Mode()); err != nil {
		_ = os.Remove(dst)

		return err
	}

	return nil
}

func copyPath(src, path, dst string, info os.FileInfo) error {
	rel, err := filepath.Rel(src, path)
	if err != nil {
		return err
	}

	newPath := filepath.Join(dst, rel)

	switch {
	case info.IsDir():
		if err := os.MkdirAll(newPath, info.Mode()); err != nil {
			return err
		}
	case info.Mode()&os.ModeSymlink != 0: // Symlink
		if rlink, err := os.Readlink(path); err != nil {
			return err
		} else if err := os.Symlink(rlink, newPath); err != nil {
			return err
		}
	default:
		if err := copyFile(path, newPath, info); err != nil {
			return err
		}
	}

	return nil
}
