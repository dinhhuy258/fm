package state

type History struct {
	loc   int
	paths []string
}

func NewHistory() *History {
	paths := make([]string, 0, 1000)

	return &History{
		loc:   -1,
		paths: paths,
	}
}

func (h *History) Push(path string) {
	if len(h.paths) == 0 || h.Peek() != path {
		h.paths = append(h.paths, path)
		h.loc++
	}
}

func (h *History) Peek() string {
	return h.paths[h.loc]
}

func (h *History) VisitLast() {
	if h.loc > 0 {
		h.loc--
	}
}

func (h *History) VisitNext() {
	if h.loc < len(h.paths)-1 {
		h.loc++
	}
}
