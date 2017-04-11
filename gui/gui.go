// Package gui ...
// to prevent linter exceptions caused by gtk librairy, theses linters are unused in the whole package: gosimple, errcheck, staticcheck, unused
// nolint: gosimple, errcheck, staticcheck, unused
package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

func buildMainWindow(win *gtk.Window) (err error) {
	win.SetTitle("Simple Example")
	label, err := gtk.LabelNew("Hello, gotk3!")
	if err != nil {
		return fmt.Errorf("unable to create label: %v", err)
	}

	// add the label to the window.
	win.Add(label)
	win.SetDefaultSize(800, 600)
	return nil
}

// GUI start the main gui window
func GUI() (err error) {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return fmt.Errorf("unable to create window: %v", err)
	}
	_, err = win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	if err != nil {
		return fmt.Errorf("unable to add destroy callback: %v", err)
	}

	if err = buildMainWindow(win); err != nil {
		return fmt.Errorf("unable to build main windows: %v", err)
	}

	win.ShowAll()

	// this block forever until main window is closed
	gtk.Main()
	return nil
}
