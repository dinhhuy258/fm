package app

type State struct {
	Main *MainState
}

type MainState struct {
	Paths         []string
	SelectedIdx   int
	NumberOfFiles int
}
