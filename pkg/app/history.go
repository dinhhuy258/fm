package app

import "github.com/dinhhuy258/fm/pkg/fs"

type History struct {
	loc     int
	entries []fs.IEntry
}

func NewHistory() *History {
	entries := make([]fs.IEntry, 0, 1000)

	return &History{
		loc:     -1,
		entries: entries,
	}
}

func (h *History) Push(entry fs.IEntry) {
	if len(h.entries) == 0 || h.Peek().GetPath() != entry.GetPath() {
		h.entries = append(h.entries, entry)
		h.loc = len(h.entries) - 1
	}
}

func (h *History) Peek() fs.IEntry {
	return h.entries[h.loc]
}

func (h *History) VisitLast() {
	if h.loc > 0 {
		h.loc--
	}
}

func (h *History) VisitNext() {
	if h.loc < len(h.entries)-1 {
		h.loc++
	}
}
