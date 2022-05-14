package controller

import set "github.com/deckarep/golang-set/v2"

type Controllers struct {
	Explorer   *ExplorerController
	Help       *HelpController
	Sellection *SelectionController
}

func CreateAllControllers() *Controllers {
	selections := set.NewSet[string]()

	return &Controllers{
		Explorer:   newExplorerController(selections),
		Help:       newHelpController(),
		Sellection: newSelectionController(selections),
	}
}
