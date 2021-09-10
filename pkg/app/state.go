package app

type State struct {
	Main       *MainState
	Selections map[string]struct{}
}

type MainState struct {
	FocusIdx      int
	NumberOfFiles int
}
