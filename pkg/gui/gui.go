package gui

import (
	"errors"
	"sync"

	"github.com/dinhhuy258/fm/pkg/fs"
	"github.com/dinhhuy258/fm/pkg/gui/view"
	"github.com/dinhhuy258/gocui"
)

type Gui struct {
	g              *gocui.Gui
	views          *view.Views
	onViewsCreated func()
}

func (gui *Gui) StartProgress(total int) {
	gui.views.Progress.StartProgress(total)
}

func (gui *Gui) UpdateProgress() {
	gui.views.Progress.AddCurrent(1)
}

func (gui *Gui) IsProgressFinished() bool {
	return gui.views.Progress.IsFinished()
}

func (gui *Gui) SetLog(log string, level view.LogLevel) {
	gui.views.Log.SetLog(log, level)
}

func (gui *Gui) SetInput(ask string, onInput func(string)) {
	gui.views.Input.SetInput(ask, func(ans string) {
		gui.views.Main.SetAsCurrentView()

		onInput(ans)
	})
}

func (gui *Gui) SetConfirmation(ask string, onConfirm func(bool)) {
	gui.views.Confirm.SetConfirmation(ask, func(ans bool) {
		gui.views.Main.SetAsCurrentView()

		onConfirm(ans)
	})
}

func (gui *Gui) RenderSelections(selections []string) {
	gui.views.Selection.RenderSelections(selections)
}

func (gui *Gui) ResetCursor() {
	_ = gui.views.Main.SetCursor(0, 0)
	_ = gui.views.Main.SetOrigin(0, 0)
}

func (gui *Gui) NextCursor() {
	_ = gui.views.Main.NextCursor()
}

func (gui *Gui) PreviousCursor() {
	_ = gui.views.Main.PreviousCursor()
}

func (gui *Gui) SetHelpTitle(title string) {
	gui.views.Help.SetTitle(title)
}

func (gui *Gui) SetHelp(keys []string, msgs []string) {
	gui.views.Help.SetHelp(keys, msgs)
}

func (gui *Gui) SetMainTitle(title string) {
	gui.views.Main.SetTitle(title)
}

func (gui *Gui) UpdateSortAndFilter() {
	gui.views.SortAndFilter.UpdateSortAndFilter()
}

func (gui *Gui) RenderEntries(entries []fs.IEntry, selections map[string]struct{}, focus int) {
	gui.views.Main.RenderEntries(entries, selections, focus)
}

var (
	gui                   *Gui
	guiInitializationOnce sync.Once
)

func InitGui(onViewsCreated func()) {
	// Make sure only one gui is created
	guiInitializationOnce.Do(func() {
		gui = &Gui{
			onViewsCreated: onViewsCreated,
		}
	})
}

func GetGui() *Gui {
	if gui == nil {
		panic("gui is not initialized")
	}

	return gui
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}

	gui.g = g
	gui.g.Cursor = true
	gui.g.InputEsc = true

	defer gui.g.Close()

	gui.g.SetManager(gocui.ManagerFunc(gui.layout))

	err = gui.g.MainLoop()

	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}

	return nil
}

func (gui *Gui) SetOnKeyFunc(onKey func(string) error) {
	gui.g.SetOnKeyFunc(onKey)
}

func (gui *Gui) Quit() error {
	return gocui.ErrQuit
}
