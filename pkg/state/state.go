package state

type State struct {
	FocusIdx      int
	NumberOfFiles int
	Selections    map[string]struct{}
	History       *History
}
