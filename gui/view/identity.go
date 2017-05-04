package view

import (
	"errors"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/api"
)

// Identity represent the login view
type Identity struct {
	Module
	builder        *gtk.Builder
	onLoginSucceed func() error
	gtkQuitOnClose bool
}

// Load load and fill all the component of the login module
func (v *Identity) Load(onLoginSucceed func() error) (err error) {
	v.builder, err = gtk.BuilderNew()
	if err != nil {
		return fmt.Errorf("unable to create builder: %v", err)
	}
	// load the view from a file
	if err = v.builder.AddFromFile("gui/view/identity.ui"); err != nil {
		return fmt.Errorf("unable to add file to builder: %v", err)
	}

	v.onLoginSucceed = onLoginSucceed
	// get window from loaded file
	v.Window, err = v.FindWindowWithBuilder(v.builder, "window_identity")
	if err != nil {
		return fmt.Errorf("unable to find window in builder: %v", err)
	}

	// attach event handlers to objects
	if err = v.attachWindowBasicSignals(); err != nil {
		return fmt.Errorf("unable to attach signals: %v", err)
	}
	if err = v.AttachButtonClickedSignal(v.builder, "button_register", v.onRegisterClicked); err != nil {
		return fmt.Errorf("unable to add button callback: %v", err)
	}
	if err = v.AttachButtonClickedSignal(v.builder, "button_login", v.onLoginClicked); err != nil {
		return fmt.Errorf("unable to add button callback: %v", err)
	}

	v.gtkQuitOnClose = true

	// finally show the window
	v.Window.ShowAll()
	return nil
}

func (v *Identity) attachWindowBasicSignals() (err error) {
	_, err = v.Window.Connect("destroy", func() error {
		if v.gtkQuitOnClose {
			gtk.MainQuit()
		}
		return nil
	}, nil)
	return err
}

func (v *Identity) onLoginClicked() (err error) {
	// get certificate from input
	fileChooserCert, err := v.FindFileChooserButtonWithBuilder(v.builder, "filechooser_certificate_login")
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to find file chooser certificate: %v", err))
	}
	cert := fileChooserCert.GetFilename()

	// get key and key password from inputs
	key, keypwd, err := v.loadKeyInputs("login")
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to load keys inputs: %v", err))
	}

	log.Debugf("selected identity key file: %q", key)
	log.Debugf("selected crt to login: %q", cert)

	// try to login
	_, err = api.API.LoginWithCertsFilename(cert, key, []byte(keypwd))
	if err != nil {
		v.Dialog(gtk.MESSAGE_ERROR, "An error occurred: %v", err)
		return err
	}

	// it's a match! hide this window and let the magic happen
	defer v.Window.Destroy()
	v.gtkQuitOnClose = false
	return log.ErrorIf(v.onLoginSucceed())
}

func (v *Identity) onRegisterClicked() (err error) {
	// get key and key password from inputs
	key, keypwd, err := v.loadKeyInputs("register")
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to load keys inputs: %v", err))
	}

	log.Debugf("selected identity key file: %q", key)

	// try to register
	_, err = api.API.RegisterWithKeyPairFilename(key, []byte(keypwd))
	if err != nil {
		v.Dialog(gtk.MESSAGE_ERROR, "An error occurred: %v", err)
		return err
	}

	defer v.Window.Destroy()
	v.gtkQuitOnClose = false
	return log.ErrorIf(v.onLoginSucceed())
}

func (v *Identity) loadKeyInputs(suffix string) (key string, keypwd string, err error) {
	fileChooserPrivKey, err := v.FindFileChooserButtonWithBuilder(v.builder, "filechooser_privkey_"+suffix)
	if err != nil {
		return "", "", fmt.Errorf("unable to find file chooser private key: %v", err)
	}
	key = fileChooserPrivKey.GetFilename()
	if key == "" {
		return "", "", errors.New("no key file selected")
	}

	entryPrivKeyPwd, err := v.FindEntryWithBuilder(v.builder, "entry_privpwd_"+suffix)
	if err != nil {
		return "", "", fmt.Errorf("unable to find entry private key password: %v", err)
	}
	keypwd, err = entryPrivKeyPwd.GetText()
	if err != nil {
		return "", "", fmt.Errorf("unable to get text from entry private key password: %v", err)
	}

	return key, keypwd, nil
}
