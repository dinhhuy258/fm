package controller

import (
	"strings"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/optional"
	"github.com/dinhhuy258/fm/pkg/pipe"
)

const markSeparator = ":"

type MarkController struct {
	*BaseController

	marks    map[string]string
	markPath string
}

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

func (mc *MarkController) SaveMark(key, path string) {
	mc.marks[key] = path
	mc.writeMarksToFile()
}

func (mc *MarkController) LoadMark(key string) optional.Optional[string] {
	path, hasKey := mc.marks[key]
	if !hasKey {
		return optional.NewEmpty[string]()
	}

	return optional.New(path)
}

func (mc *MarkController) writeMarksToFile() {
	lines := make([]string, 0)

	for k, v := range mc.marks {
		lines = append(lines, k+markSeparator+v)
	}

	fs.WriteToFile(mc.markPath, lines, true)
}
