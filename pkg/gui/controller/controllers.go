package controller

type Controllers struct {
	Explorer *ExplorerController
	Help     *HelpController
}

func CreateAllControllers() *Controllers {
	return &Controllers{
		Explorer: newExplorerController(),
		Help:     newHelpController(),
	}
}
