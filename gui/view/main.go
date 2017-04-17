// Package view ...
// to prevent linter exceptions caused by gtk librairy, theses linters are unused in the whole package: gosimple, errcheck, staticcheck, unused
// nolint: gosimple, errcheck, staticcheck, unused
package view

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

// Main represent the main view
type Main struct {
	Module
	builder *gtk.Builder
}

// Load load and fill all the component of the main module
func (v *Main) Load() (err error) {
	v.builder, err = gtk.BuilderNew()
	if err != nil {
		return fmt.Errorf("unable to create builder: %v", err)
	}

	if err = v.builder.AddFromFile("gui/view/main.ui"); err != nil {
		return fmt.Errorf("unable to add file to builder: %v", err)
	}

	v.Window, err = v.FindWindowWithBuilder(v.builder, "window_main")
	if err != nil {
		return fmt.Errorf("unable to find window in builder: %v", err)
	}

	if err = v.attachSignals(); err != nil {
		return fmt.Errorf("unable to attach signals: %v", err)
	}

	v.Window.ShowAll()
	return nil
}

func (v *Main) attachSignals() (err error) {
	_, err = v.Window.Connect("destroy", func() { gtk.MainQuit() }, nil)
	if err != nil {
		return fmt.Errorf("unable to attach destroy signal to window: %v", err)
	}
	return err
}
