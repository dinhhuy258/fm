package app

type State struct {
	FocusIdx      int
	//TODO: Remove this variable, I don't think we need to keep track of this variable here
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
