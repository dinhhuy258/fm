package fs

import (
	"fmt"
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

func Dir(path string) string {
	return filepath.Dir(path)
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

func WriteToFile(filePath string, lines []string, override bool) {
	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	if override {
		flags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	}

	file, err := os.OpenFile(filePath, flags, 0o644)
	if err != nil {
		return
	}

	defer file.Close()

	for _, line := range lines {
		if _, err = file.WriteString(line + "\n"); err != nil {
			return
		}
	}
}

func isHidden(filename string) bool {
	return filename[0:1] == "."
}