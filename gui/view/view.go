package view

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"
)

// Module represent a module of the application (login, chat, ...)
type Module struct {
	WindowBaseTitle string
	Window          *gtk.Window
}

// OnClickEvent is the prototype of a button clicked event
type OnClickEvent func() error

// Dialog open a modal box and display a message
func (m *Module) Dialog(messageType gtk.MessageType, format string, args ...interface{}) {
	switch messageType {
	case gtk.MESSAGE_WARNING:
		log.Warningf(format, args...)
	case gtk.MESSAGE_INFO:
		log.Infof(format, args...)
	case gtk.MESSAGE_ERROR:
		log.Errorf(format, args...)
	}
	infoBox := gtk.MessageDialogNew(m.Window, gtk.DIALOG_DESTROY_WITH_PARENT, messageType, gtk.BUTTONS_CLOSE, format, args...)
	_, err := infoBox.Connect("response", infoBox.Destroy)
	if err != nil {
		log.Warningf("unable to connect response event: %v", err)
	}
	infoBox.Show()
}

// FindButtonWithBuilder return a button stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindButtonWithBuilder(builder *gtk.Builder, buttonName string) (button *gtk.Button, err error) {
	widget, err := builder.GetObject(buttonName)
	if err != nil {
		return nil, fmt.Errorf("unable to get button %q from builder: %v", buttonName, err)
	}

	button, ok := widget.(*gtk.Button)
	if !ok {
		return nil, fmt.Errorf("unable to cast button from widget")
	}

	return button, nil
}

// FindEntryWithBuilder return an entry stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindEntryWithBuilder(builder *gtk.Builder, entryName string) (button *gtk.Entry, err error) {
	widget, err := builder.GetObject(entryName)
	if err != nil {
		return nil, fmt.Errorf("unable to get file chooser %q from builder: %v", entryName, err)
	}

	entry, ok := widget.(*gtk.Entry)
	if !ok {
		return nil, fmt.Errorf("unable to cast file chooser from widget")
	}

	return entry, nil
}

// FindFileChooserButtonWithBuilder return an file chooser button stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindFileChooserButtonWithBuilder(builder *gtk.Builder, fileChooserName string) (button *gtk.FileChooserButton, err error) {
	widget, err := builder.GetObject(fileChooserName)
	if err != nil {
		return nil, fmt.Errorf("unable to get file chooser %q from builder: %v", fileChooserName, err)
	}

	button, ok := widget.(*gtk.FileChooserButton)
	if !ok {
		return nil, fmt.Errorf("unable to cast file chooser from widget")
	}

	return button, nil
}

// FindWindowWithBuilder return an window stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindWindowWithBuilder(builder *gtk.Builder, windowName string) (window *gtk.Window, err error) {
	widget, err := builder.GetObject(windowName)
	if err != nil {
		return nil, fmt.Errorf("unable to get window %q from builder: %v", windowName, err)
	}

	window, ok := widget.(*gtk.Window)
	if !ok {
		return nil, fmt.Errorf("unable to cast window from widget")
	}

	return window, nil
}
