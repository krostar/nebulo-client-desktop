package view

import (
	"fmt"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/user"
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

	if err = v.attachWindowBasicSignals(); err != nil {
		return fmt.Errorf("unable to attach signals: %v", err)
	}

	if err = v.attachMenuSignals(); err != nil {
		return fmt.Errorf("unable to attach menu signals: %v", err)
	}

	v.Window.ShowAll()
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
