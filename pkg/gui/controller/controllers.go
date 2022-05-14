package controller

import set "github.com/deckarep/golang-set/v2"

type Controllers struct {
	Explorer   *ExplorerController
	Help       *HelpController
	Sellection *SelectionController
	Progress   *ProgressController
	Log        *LogController
	Input      *InputController
}

func CreateAllControllers() *Controllers {
	// Selections object to share between explorer and selection controllers
	selections := set.NewSet[string]()

	return &Controllers{
		Explorer:   newExplorerController(selections),
		Sellection: newSelectionController(selections),
		Help:       newHelpController(),
		Progress:   newProgressController(),
		Log:        newLogController(),
		Input:      newInputController(),
	}
}
