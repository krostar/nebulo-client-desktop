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
	// load the view from a file
	if err = v.builder.AddFromFile("gui/view/main.ui"); err != nil {
		return fmt.Errorf("unable to add file to builder: %v", err)
	}

	// get window from loaded file
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
