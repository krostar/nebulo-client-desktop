package view

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/api"
	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/user"
)

// ChannelAdd represent the channel addition view
type ChannelAdd struct {
	Module
	builder   *gtk.Builder
	dialog    *gtk.Dialog
	treeview  *gtk.TreeView
	liststore *gtk.ListStore
}

// Load load and fill all the component of the add channel module
func (v *ChannelAdd) Load(parent *gtk.Window) (err error) {
	v.builder, err = gtk.BuilderNew()
	if err != nil {
		return fmt.Errorf("unable to create builder: %v", err)
	}
	// load the view from a file
	if err = v.builder.AddFromFile("gui/view/channel_add.ui"); err != nil {
		return fmt.Errorf("unable to add file to builder: %v", err)
	}

	// get dialog from loaded file
	v.dialog, err = v.FindDialogWithBuilder(v.builder, "dialog_channel")
	if err != nil {
		return fmt.Errorf("unable to find dialog in builder: %v", err)
	}
	v.dialog.SetTitle(v.WindowBaseTitle + "Channel add")
	v.dialog.SetTransientFor(parent)

	if err = v.AttachButtonClickedSignal(v.builder, "button_cancel", v.onCancelClicked); err != nil {
		return fmt.Errorf("unable to attach signals: %v", err)
	}
	if err = v.AttachButtonClickedSignal(v.builder, "button_add", v.onAddClicked); err != nil {
		return fmt.Errorf("unable to attach signals: %v", err)
	}

	v.treeview, err = v.FindTreeViewWithBuilder(v.builder, "treeview_contacts")
	if err != nil {
		return fmt.Errorf("unable to find treeview in builder: %v", err)
	}

	err = v.fillContacts(config.Config.Run.ContactsFile)
	if err != nil {
		return fmt.Errorf("unable to fill treeview with contacts: %v", err)
	}

	v.dialog.Show()
	return nil
}

func (v *ChannelAdd) fillContacts(file string) (err error) {
	v.liststore, err = gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		return fmt.Errorf("unable to create list store: %v", err)
	}

	for _, contact := range user.Logged.Contacts {
		iter := v.liststore.Append()
		err = v.liststore.Set(iter, []int{0}, []interface{}{contact.Name})
		if err != nil {
			return fmt.Errorf("unable to insert contact %q: %v", contact.Name, err)
		}
	}
	v.treeview.SetModel(v.liststore)
	return nil
}

func (v *ChannelAdd) onCancelClicked() (err error) {
	v.dialog.Destroy()
	return nil
}

func (v *ChannelAdd) onAddClicked() (err error) {
	entryChannelName, err := v.FindEntryWithBuilder(v.builder, "entry_channel_name")
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to find entry channel name: %v", err))
	}
	channelName, err := entryChannelName.GetText()
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to get text from entry channel name: %v", err))
	}

	selection, err := v.treeview.GetSelection()
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to get selection from tree view: %v", err))
	}
	selected := selection.GetSelectedRows(v.liststore)

	var channelMembersPkey []string
	selected.Foreach(func(item interface{}) {
		treepath, ok := item.(*gtk.TreePath)
		if !ok {
			return
		}
		channelMembersPkey = append(channelMembersPkey, user.Logged.Contacts[treepath.GetIndices()[0]].PublicKeyB64)
	})

	_, err = api.API.ChannelCreate(channelName, channelMembersPkey)
	if err != nil {
		v.Dialog(gtk.MESSAGE_ERROR, "An error occurred: %v", err)
		return err
	}

	v.dialog.Destroy()
	return nil
}
