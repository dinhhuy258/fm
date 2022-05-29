package controller

import (
	"strings"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/optional"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

const markSeparator = ":"

// MarkController is the controller for mark
type MarkController struct {
	*BaseController

	marks    map[string]string
	markPath string
}

// newMarkController creates a new mark controller
func newMarkController(baseController *BaseController, pipe *pipe.Pipe) *MarkController {
	markPath := pipe.GetMarkPath()
	markLines := fs.ReadFromFile(markPath)
	marks := make(map[string]string)

	// Load the marks from the file
	for _, line := range markLines {
		if line == "" {
			continue
		}

		// The mark line is in the format of: <key>:<path>
		toks := strings.SplitN(line, markSeparator, 2)
		if len(toks) != 2 {
			// Invalid mark line
			continue
		}

		marks[toks[0]] = toks[1]
	}

	return &MarkController{
		BaseController: baseController,

		markPath: markPath,
		marks:    marks,
	}
}

// SaveMark saves the mark to mark map as well as the mark file
func (mc *MarkController) SaveMark(key, path string) {
	mc.marks[key] = path
	mc.writeMarksToFile()
}

// LoadMark loads the mark from the mark map
func (mc *MarkController) LoadMark(key string) optional.Optional[string] {
	path, hasKey := mc.marks[key]
	if !hasKey {
		return optional.NewEmpty[string]()
	}

	return optional.New(path)
}

// writeMarksToFile writes the marks to the mark file
func (mc *MarkController) writeMarksToFile() {
	lines := make([]string, 0)

	for k, v := range mc.marks {
		lines = append(lines, k+markSeparator+v)
	}

	fs.WriteToFile(mc.markPath, lines, true)
}
