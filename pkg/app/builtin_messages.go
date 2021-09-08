package app

func focusNext(app *App) error {
	if app.State.Main.SelectedIdx == app.State.Main.NumberOfFiles {
		return nil
	}

	v := app.Gui.Views.Main

	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}

	app.State.Main.SelectedIdx++

	return nil
}

func focusPrevious(app *App) error {
	if app.State.Main.SelectedIdx == 1 {
		return nil
	}

	v := app.Gui.Views.Main

	ox, oy := v.Origin()
	cx, cy := v.Cursor()

	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}

	app.State.Main.SelectedIdx--

	return nil
}

func enter(app *App) error {
	currentNode := app.FileManager.Dir.Nodes[app.State.Main.SelectedIdx-1]

	if currentNode.IsDir {
		app.FileManager.LoadDirectory(currentNode.AbsolutePath)
	}

	return nil
}

func back(app *App) error {
	parent := app.FileManager.Dir.Parent()

	app.FileManager.LoadDirectory(parent)

	return nil
}
