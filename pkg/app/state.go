package app

type State struct {
	FocusIdx      int
	NumberOfFiles int
	Selections    map[string]struct{}
}
