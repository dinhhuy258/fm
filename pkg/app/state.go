package app

type State struct {
	FocusIdx      int
	NumberOfFiles int
	Selections    map[string]struct{}
	History       *History
	Marks         map[string]string
}

func NewState() *State {
	return &State{
		FocusIdx:      -1,
		NumberOfFiles: 0,
		Selections:    map[string]struct{}{},
		Marks:         map[string]string{},
		History:       NewHistory(),
	}
}
