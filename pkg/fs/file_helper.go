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

func CreateFile(name string) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
