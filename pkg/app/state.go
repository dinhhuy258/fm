package app

type State struct {
	Main *MainState
}

type MainState struct {
	FocusIdx      int
	NumberOfFiles int
}
