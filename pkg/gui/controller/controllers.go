package controller

type Controllers struct {
	Explorer *ExplorerController
}

func CreateAllControllers() *Controllers {
	return &Controllers{
		Explorer: newExplorerController(),
	}
}
