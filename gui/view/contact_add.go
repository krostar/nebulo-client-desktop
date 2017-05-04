package view

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/contact"
)

// ContactAdd represent the contact addition view
type ContactAdd struct {
	Module
	builder *gtk.Builder
	dialog  *gtk.Dialog
}

// Load load and fill all the component of the add contact module
func (v *ContactAdd) Load(parent *gtk.Window) (err error) {
	v.builder, err = gtk.BuilderNew()
	if err != nil {
		return fmt.Errorf("unable to create builder: %v", err)
	}
	// load the view from a file
	if err = v.builder.AddFromFile("gui/view/contact_add.ui"); err != nil {
		return fmt.Errorf("unable to add file to builder: %v", err)
	}

	// get dialog from loaded file
	v.dialog, err = v.FindDialogWithBuilder(v.builder, "dialog_contact")
	if err != nil {
		return fmt.Errorf("unable to find dialog in builder: %v", err)
	}
	v.dialog.SetTitle(v.WindowBaseTitle + "Contact add")
	v.dialog.SetTransientFor(parent)

	if err = v.AttachButtonClickedSignal(v.builder, "button_cancel", v.onCancelClicked); err != nil {
		return fmt.Errorf("unable to attach signals: %v", err)
	}
	if err = v.AttachButtonClickedSignal(v.builder, "button_add", v.onAddClicked); err != nil {
		return fmt.Errorf("unable to attach signals: %v", err)
	}

	v.dialog.Show()
	return nil
}

func (v *ContactAdd) onCancelClicked() (err error) {
	v.dialog.Destroy()
	return nil
}

func (v *ContactAdd) onAddClicked() (err error) {
	entryContactName, err := v.FindEntryWithBuilder(v.builder, "entry_display_name")
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to find entry display name: %v", err))
	}
	contactName, err := entryContactName.GetText()
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to get text from entry display name: %v", err))
	}

	textviewContactPK, err := v.FindTextViewWithBuilder(v.builder, "textview_publickey_b64")
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to find textview public key: %v", err))
	}

	tvContactPKBuffer, err := textviewContactPK.GetBuffer()
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to get text from textview public key: %v", err))
	}

	tvItersStart, tvItersStop := tvContactPKBuffer.GetBounds()
	contactPK, err := tvContactPKBuffer.GetText(tvItersStart, tvItersStop, false)
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to get text from text view contact public key: %v", err))
	}

	log.Debugf("new contact: %q, %q", contactName, contactPK)
	_, err = contact.AddToFile(config.Config.Run.ContactsFile, contactName, contactPK)
	if err != nil {
		v.Dialog(gtk.MESSAGE_ERROR, "An error occurred: %v", err)
		return err
	}

	v.dialog.Destroy()
	return nil
}
