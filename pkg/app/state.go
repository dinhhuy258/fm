package app

type State struct {
	Main *MainState
}

type MainState struct {
	SelectedIdx   int
	NumberOfFiles int
}
