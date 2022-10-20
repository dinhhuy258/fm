package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// Humanize returns a human-readable string of the size.
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

// CreateFile creates a file with the given path,
// if the override flag is true, the file will be overwritten.
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

// Dir returns the parent directory of the given path.
func Dir(path string) string {
	return filepath.Dir(path)
}

// IsDir returns true if the given path is a directory.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// WriteToFile writes the given lines to the given file.
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

func IsPathExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}
