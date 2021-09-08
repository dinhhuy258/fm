package app

type History struct {
	loc   int
	paths []string
}

func NewHistory(initializePath string) *History {
	paths := make([]string, 0, 1000)
	paths = append(paths, initializePath)

	return &History{
		loc:   0,
		paths: paths,
	}
}

func (h *History) Push(path string) {
	if h.Peek() != path {
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
