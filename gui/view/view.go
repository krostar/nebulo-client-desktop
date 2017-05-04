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

// AttachButtonClickedSignal attach the clicked sign to a button
func (m *Module) AttachButtonClickedSignal(builder *gtk.Builder, buttonName string, onClick OnClickEvent) (err error) {
	button, err := m.FindButtonWithBuilder(builder, buttonName)
	if err != nil {
		return fmt.Errorf("unable to find button %q: %v", buttonName, err)
	}

	if _, err = button.Connect("clicked", onClick); err != nil {
		return fmt.Errorf("unable to add signal clicked to button %q: %v", buttonName, err)
	}

	return nil
}

// FindTextViewWithBuilder return a text view stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindTextViewWithBuilder(builder *gtk.Builder, tvName string) (tv *gtk.TextView, err error) {
	widget, err := builder.GetObject(tvName)
	if err != nil {
		return nil, fmt.Errorf("unable to get text view %q from builder: %v", tvName, err)
	}

	tv, ok := widget.(*gtk.TextView)
	if !ok {
		return nil, fmt.Errorf("unable to cast text view from widget")
	}

	return tv, nil
}

// FindListBoxWithBuilder return a list box stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindListBoxWithBuilder(builder *gtk.Builder, tvName string) (tv *gtk.ListBox, err error) {
	widget, err := builder.GetObject(tvName)
	if err != nil {
		return nil, fmt.Errorf("unable to get list box %q from builder: %v", tvName, err)
	}

	tv, ok := widget.(*gtk.ListBox)
	if !ok {
		return nil, fmt.Errorf("unable to cast list box from widget")
	}

	return tv, nil
}

// FindListStoreWithBuilder return a list store stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindListStoreWithBuilder(builder *gtk.Builder, tvName string) (tv *gtk.ListStore, err error) {
	widget, err := builder.GetObject(tvName)
	if err != nil {
		return nil, fmt.Errorf("unable to get list store %q from builder: %v", tvName, err)
	}

	tv, ok := widget.(*gtk.ListStore)
	if !ok {
		return nil, fmt.Errorf("unable to cast list store from widget")
	}

	return tv, nil
}

// FindTreeViewWithBuilder return a tree view stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindTreeViewWithBuilder(builder *gtk.Builder, tvName string) (tv *gtk.TreeView, err error) {
	widget, err := builder.GetObject(tvName)
	if err != nil {
		return nil, fmt.Errorf("unable to get tree view %q from builder: %v", tvName, err)
	}

	tv, ok := widget.(*gtk.TreeView)
	if !ok {
		return nil, fmt.Errorf("unable to cast tree view from widget")
	}

	return tv, nil
}

// FindPanedWithBuilder return a paned stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindPanedWithBuilder(builder *gtk.Builder, panedName string) (menuItem *gtk.Paned, err error) {
	widget, err := builder.GetObject(panedName)
	if err != nil {
		return nil, fmt.Errorf("unable to get paned %q from builder: %v", panedName, err)
	}

	paned, ok := widget.(*gtk.Paned)
	if !ok {
		return nil, fmt.Errorf("unable to cast paned from widget")
	}

	return paned, nil
}

// FindMenuItemWithBuilder return a menuitem stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindMenuItemWithBuilder(builder *gtk.Builder, imiName string) (menuItem *gtk.MenuItem, err error) {
	widget, err := builder.GetObject(imiName)
	if err != nil {
		return nil, fmt.Errorf("unable to get menu item %q from builder: %v", imiName, err)
	}

	menuItem, ok := widget.(*gtk.MenuItem)
	if !ok {
		return nil, fmt.Errorf("unable to cast menuItem from widget")
	}

	return menuItem, nil
}

// FindDialogWithBuilder return a dialog stored in a builder, based on his name
// nolint: dupl
func (m *Module) FindDialogWithBuilder(builder *gtk.Builder, dialogName string) (dialog *gtk.Dialog, err error) {
	widget, err := builder.GetObject(dialogName)
	if err != nil {
		return nil, fmt.Errorf("unable to get dialog %q from builder: %v", dialogName, err)
	}

	dialog, ok := widget.(*gtk.Dialog)
	if !ok {
		return nil, fmt.Errorf("unable to cast dialog from widget")
	}

	return dialog, nil
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
