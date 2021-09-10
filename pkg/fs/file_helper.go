package fs

import "fmt"

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
