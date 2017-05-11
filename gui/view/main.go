package view

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"
	"github.com/krostar/nebulo-golib/tools/cert"
	"github.com/krostar/nebulo-golib/tools/crypto"

	"github.com/krostar/nebulo-client-desktop/api"
	"github.com/krostar/nebulo-client-desktop/channel"
	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/message"
	"github.com/krostar/nebulo-client-desktop/user"
)

// Main represent the main view
type Main struct {
	Module
	builder           *gtk.Builder
	channelsTreeview  *gtk.TreeView
	channelsListstore *gtk.ListStore
	messagesTreeview  *gtk.TreeView
	messagesListstore *gtk.ListStore
	messageEntry      *gtk.Entry
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

	if err = v.attachWindowBasicSignals(); err != nil {
		return fmt.Errorf("unable to attach signals: %v", err)
	}

	if err = v.attachMenuSignals(); err != nil {
		return fmt.Errorf("unable to attach menu signals: %v", err)
	}

	if err = v.createMessageList(); err != nil {
		return fmt.Errorf("unable to create channels list: %v", err)
	}

	if err = v.createChannelList(); err != nil {
		return fmt.Errorf("unable to create channels list: %v", err)
	}

	v.messageEntry, err = v.FindEntryWithBuilder(v.builder, "entry_message")
	if err != nil {
		return fmt.Errorf("unable to find button in builder: %v", err)
	}
	if err = v.makeMessageEntryUneditable(); err != nil {
		return fmt.Errorf("unable to make message entry uneditable: %v", err)
	}
	if _, err = v.messageEntry.Connect("activate", v.onMessageSent, nil); err != nil {
		return fmt.Errorf("unable to connect signal activate to message entry: %v", err)
	}

	v.Window.ShowAll()
	return nil
}

func (v *Main) makeMessageEntryUneditable() (err error) {
	if err = v.messageEntry.SetProperty("editable", false); err != nil {
		return fmt.Errorf("unable to make message entry uneditable")
	}
	v.messageEntry.SetSensitive(false)
	v.messageEntry.SetText("Pleace, select a channel")
	return nil
}

func (v *Main) makeMessageEntryEditable() (err error) {
	if err = v.messageEntry.SetProperty("editable", true); err != nil {
		return fmt.Errorf("unable to make message entry editable")
	}
	v.messageEntry.SetSensitive(true)
	v.messageEntry.SetText("")
	return nil
}

func (v *Main) onMessageSent() (err error) {
	entryMessage, err := v.FindEntryWithBuilder(v.builder, "entry_message")
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to find entry message: %v", err))
	}
	msg, err := entryMessage.GetText()
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to get text from entry message: %v", err))
	}
	entryMessage.SetText("")

	selection, err := v.channelsTreeview.GetSelection()
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to get selection from treeview: %v", err))
	}
	channelName, err := v.getChannelFromSelectedChannelList(selection)
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to get selected channel name: %v", err))
	}

	log.Debugf("Channel: %q -- Message: %q", channelName, msg)
	if err = api.API.MessageCreate(channelName, msg); err != nil {
		return log.ErrorIf(fmt.Errorf("unable to send message to server: %v", err))
	}
	return nil
}

func (v *Main) onChannelSelectionChanged(selection *gtk.TreeSelection) (err error) {
	channelName, err := v.getChannelFromSelectedChannelList(selection)
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to get selected channel name: %v", err))
	}

	description, err := v.FindLabelWithBuilder(v.builder, "description_conv")
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to find channel description label: %v", err))
	}
	description.SetText(fmt.Sprintf("Select channel: %s", channelName))

	messages, err := api.API.MessageList(channelName, time.Time{})
	if err != nil {
		v.Dialog(gtk.MESSAGE_ERROR, "An error occured on messages fetching: %v", err)
		return log.ErrorIf(err)
	}
	if err = v.MessagesRefresh(messages); err != nil {
		return log.ErrorIf(fmt.Errorf("unable to refresh messages: %v", err))
	}

	if err = v.makeMessageEntryEditable(); err != nil {
		return log.ErrorIf(fmt.Errorf("unable to make message entry editable: %v", err))
	}
	log.Debugf("new channel selected: %v", channelName)
	return nil
}

func (v *Main) ChannelsRefresh() (err error) {
	v.channelsListstore.Clear()

	for cName := range channel.Channels {
		iter := v.channelsListstore.Append()
		err = v.channelsListstore.Set(iter, []int{0}, []interface{}{cName})
		if err != nil {
			return fmt.Errorf("unable to insert channel %q: %v", cName, err)
		}
	}
	return nil
}

func (v *Main) MessagesRefresh(messages []*message.Message) (err error) {
	v.messagesListstore.Clear()
	pKeyPem, err := cert.ParsePrivateKeyPEMFromFile(config.Config.Run.TLS.Key, []byte(config.Config.Run.TLS.KeyPassword))
	if err != nil {
		return fmt.Errorf("unable to decode PEM encoded private key file %q: %v", config.Config.Run.TLS.Key, err)
	}
	pKey, ok := pKeyPem.(*rsa.PrivateKey)
	if !ok {
		return errors.New("cant cast private key to rsa private key")
	}

	for _, m := range messages {
		log.Debugln(m.Ciphertext, m.Keys, m.Integrity)
		plaintext, err := crypto.Decrypt(m.Ciphertext, m.Keys, m.Integrity, *pKey)
		if err != nil {
			return fmt.Errorf("unable to decrypt: %v", err)
		}
		m.Plaintext = string(plaintext)
		iter := v.messagesListstore.Append()
		err = v.messagesListstore.Set(iter, []int{0}, []interface{}{
			fmt.Sprintf("%s (%s): %s", m.Sender.KeyFingerprint, m.Sender.DisplayName, m.Plaintext),
		})
		if err != nil {
			return fmt.Errorf("unable to insert channel %q (from %q): %v", m.Plaintext, m.Sender.DisplayName, err)
		}
	}
	return nil
}

func (v *Main) getChannelFromSelectedChannelList(selection *gtk.TreeSelection) (channelName string, err error) {
	model, iter, ok := selection.GetSelected()
	if !ok {
		return "", errors.New("ok is false on channel selection")
	}
	ivalue, err := model.(*gtk.TreeModel).GetValue(iter, 0)
	if err != nil {
		return "", fmt.Errorf("unable to get channel name value from tree model: %v", err)
	}
	channelName, err = ivalue.GetString()
	if err != nil {
		return "", fmt.Errorf("unable to get channel name to string: %v", err)
	}
	return channelName, nil
}

func (v *Main) createChannelList() (err error) {
	v.channelsTreeview, err = v.FindTreeViewWithBuilder(v.builder, "treeview_channels")
	if err != nil {
		return fmt.Errorf("unable to find listbox channel: %v", err)
	}
	v.channelsListstore, err = gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		return fmt.Errorf("unable to create list store: %v", err)
	}
	v.channelsTreeview.SetModel(v.channelsListstore)

	selection, err := v.channelsTreeview.GetSelection()
	if err != nil {
		return fmt.Errorf("unable to get selection from treeview: %v", err)
	}
	if _, err = selection.Connect("changed", v.onChannelSelectionChanged, nil); err != nil {
		return fmt.Errorf("unable to attach destroy signal to channels treeview: %v", err)
	}

	return nil
}

func (v *Main) createMessageList() (err error) {
	v.messagesTreeview, err = v.FindTreeViewWithBuilder(v.builder, "treeview_messages")
	if err != nil {
		return fmt.Errorf("unable to find listbox message: %v", err)
	}
	v.messagesListstore, err = gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		return fmt.Errorf("unable to create list store: %v", err)
	}
	v.messagesTreeview.SetModel(v.messagesListstore)
	return nil
}

func (v *Main) attachWindowBasicSignals() (err error) {
	if _, err = v.Window.Connect("destroy", func() error { gtk.MainQuit(); return nil }, nil); err != nil {
		return fmt.Errorf("unable to attach destroy signal to window: %v", err)
	}

	paned, err := v.FindPanedWithBuilder(v.builder, "paned_main")
	if err != nil {
		return fmt.Errorf("unable to find app quit menu item: %v", err)
	}
	if _, err = paned.Connect("accept-position", func() error {
		log.Debugln("handled moved")
		return nil
	}, nil); err != nil {
		return fmt.Errorf("unable to attach destroy signal to window: %v", err)
	}
	return err
}

func (v *Main) attachMenuSignals() (err error) {
	if err = v.attachMenuAppSignals(); err != nil {
		return err
	}
	if err = v.attachMenuProfilSignals(); err != nil {
		return err
	}
	if err = v.attachMenuChannelsSignals(); err != nil {
		return err
	}
	if err = v.attachMenuContactsSignals(); err != nil {
		return err
	}
	if err = v.attachMenuSettingsSignals(); err != nil {
		return err
	}
	return nil
}

func (v *Main) attachMenuAppSignals() (err error) {
	appQuit, err := v.FindMenuItemWithBuilder(v.builder, "menuitem_app_quit")
	if err != nil {
		return fmt.Errorf("unable to find app quit menu item: %v", err)
	}
	if _, err = appQuit.Connect("activate", func() error {
		v.Window.Destroy()
		return nil
	}, nil); err != nil {
		return fmt.Errorf("unable to attach activate signal to app quit menuitem: %v", err)
	}
	return nil
}

func (v *Main) attachMenuProfilSignals() (err error) {
	profilSee, err := v.FindMenuItemWithBuilder(v.builder, "menuitem_profil_see")
	if err != nil {
		return fmt.Errorf("unable to find profil see menu item: %v", err)
	}
	profilEdit, err := v.FindMenuItemWithBuilder(v.builder, "menuitem_profil_edit")
	if err != nil {
		return fmt.Errorf("unable to find profil edit menu item: %v", err)
	}
	profilShare, err := v.FindMenuItemWithBuilder(v.builder, "menuitem_profil_share")
	if err != nil {
		return fmt.Errorf("unable to find profil share menu item: %v", err)
	}

	if _, err = profilSee.Connect("activate", func() error {
		return nil
	}, nil); err != nil {
		return fmt.Errorf("unable to attach activate signal to profil see menuitem: %v", err)
	}
	if _, err = profilEdit.Connect("activate", func() error {
		return nil
	}, nil); err != nil {
		return fmt.Errorf("unable to attach activate signal to profil edit menuitem: %v", err)
	}
	if _, err = profilShare.Connect("activate", func() error {
		var clipboard *gtk.Clipboard
		clipboard, err = gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
		if err != nil {
			return log.ErrorIf(fmt.Errorf("unable to get clipboard: %v", err))
		}
		clipboard.SetText(user.Logged.PublicKeyDerBase64)
		return nil
	}, nil); err != nil {
		return fmt.Errorf("unable to attach activate signal to profil share menuitem: %v", err)
	}
	return nil
}

func (v *Main) attachMenuChannelsSignals() (err error) {
	channelsCreate, err := v.FindMenuItemWithBuilder(v.builder, "menuitem_channels_create")
	if err != nil {
		return fmt.Errorf("unable to find channels create menu item: %v", err)
	}
	if _, err = channelsCreate.Connect("activate", func() error {
		channelDialog := &ChannelAdd{Module: v.Module}
		return log.ErrorIf(channelDialog.Load(v.Window))
	}, nil); err != nil {
		return fmt.Errorf("unable to attach activate signal to channels create menuitem: %v", err)
	}
	return nil
}

func (v *Main) attachMenuContactsSignals() (err error) {
	contactsAdd, err := v.FindMenuItemWithBuilder(v.builder, "menuitem_contacts_add")
	if err != nil {
		return fmt.Errorf("unable to find contacts add menu item: %v", err)
	}

	if _, err = contactsAdd.Connect("activate", func() error {
		contactDialog := &ContactAdd{Module: v.Module}
		return log.ErrorIf(contactDialog.Load(v.Window))
	}, nil); err != nil {
		return fmt.Errorf("unable to attach activate signal to contacts add menuitem: %v", err)
	}
	return nil
}

func (v *Main) attachMenuSettingsSignals() (err error) {
	settingsHelp, err := v.FindMenuItemWithBuilder(v.builder, "menuitem_settings_help")
	if err != nil {
		return fmt.Errorf("unable to find settings help menu item: %v", err)
	}
	if _, err = settingsHelp.Connect("activate", func() error {
		return nil
	}, nil); err != nil {
		return fmt.Errorf("unable to attach activate signal to settings help menuitem: %v", err)
	}
	return nil
}
