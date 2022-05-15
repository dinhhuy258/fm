package command

func MarkSave(app IApp, params ...interface{}) error {
	appGui := app.GetGui()
	explorerController := appGui.GetControllers().Explorer

	key, _ := params[0].(string)
	// Exit mark mode
	_ = app.PopMode()
	entry := explorerController.GetCurrentEntry()
	app.MarkSave(key, entry.GetPath())

	return nil
}

func MarkLoad(app IApp, params ...interface{}) error {
	key, _ := params[0].(string)
	// Exit mark mode
	_ = app.PopMode()

	if path, hasKey := app.MarkLoad(key); hasKey {
		return FocusPath(app, path)
	}

	return nil
}
