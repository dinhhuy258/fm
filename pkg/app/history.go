package app

import "github.com/dinhhuy258/fm/pkg/fs"

type History struct {
	loc   int
	nodes []*fs.Node
}

func NewHistory() *History {
	nodes := make([]*fs.Node, 0, 1000)

	return &History{
		loc:   -1,
		nodes: nodes,
	}
}

func (h *History) Push(node *fs.Node) {
	if len(h.nodes) == 0 || h.Peek().AbsolutePath != node.AbsolutePath {
		h.nodes = append(h.nodes, node)
		h.loc = len(h.nodes) - 1
	}
}

func (h *History) Peek() *fs.Node {
	return h.nodes[h.loc]
}

func (h *History) VisitLast() {
	if h.loc > 0 {
		h.loc--
	}
}

func (h *History) VisitNext() {
	if h.loc < len(h.nodes)-1 {
		h.loc++
	}
}
