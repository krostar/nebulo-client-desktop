// Package view ...
// to prevent linter exceptions caused by gtk librairy, theses linters are unused in the whole package: gosimple, errcheck, staticcheck, unused
// nolint: gosimple, errcheck, staticcheck, unused
package view

import (
	"errors"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/api"
	"github.com/krostar/nebulo-client-desktop/config"
)

// Login represent the login view
type Login struct {
	Module
	builder *gtk.Builder
}

// Load load and fill all the component of the login module
func (v *Login) Load() (err error) {
	v.builder, err = gtk.BuilderNew()
	if err != nil {
		return fmt.Errorf("unable to create builder: %v", err)
	}

	if err = v.builder.AddFromFile("gui/view/login.ui"); err != nil {
		return fmt.Errorf("unable to add file to builder: %v", err)
	}

	v.Window, err = v.FindWindowWithBuilder(v.builder, "window_identity")
	if err != nil {
		return fmt.Errorf("unable to find window in builder: %v", err)
	}

	if err = v.attachWindowBasicSignals(); err != nil {
		return fmt.Errorf("unable to attach signals: %v", err)
	}

	if err = v.attachButtonClickedSignals(v.builder, "button_register", v.onRegisterClicked); err != nil {
		return fmt.Errorf("unable to add button callback: %v", err)
	}

	if err = v.attachButtonClickedSignals(v.builder, "button_login", v.onLoginClicked); err != nil {
		return fmt.Errorf("unable to add button callback: %v", err)
	}

	v.Window.ShowAll()
	return nil
}

func (v *Login) attachWindowBasicSignals() (err error) {
	_, err = v.Window.Connect("destroy", func() { gtk.MainQuit() }, nil)
	if err != nil {
		return fmt.Errorf("unable to attach destroy signal to window: %v", err)
	}
	return err
}

func (v *Login) attachButtonClickedSignals(builder *gtk.Builder, buttonName string, onClick OnClickEvent) (err error) {
	button, err := v.FindButtonWithBuilder(builder, buttonName)
	if err != nil {
		return fmt.Errorf("unable to find button %q: %v", buttonName, err)
	}

	if _, err = button.Connect("clicked", onClick); err != nil {
		return fmt.Errorf("unable to add signal clicked to button %q: %v", buttonName, err)
	}

	return nil
}

func (v *Login) onLoginClicked() (err error) {
	var (
		cert   string
		key    = config.Config.Run.TLS.Key
		keypwd = config.Config.Run.TLS.KeyPassword
	)

	fileChooserCert, err := v.FindFileChooserButtonWithBuilder(v.builder, "filechooser_certificate")
	if err != nil {
		err = fmt.Errorf("unable to find file chooser certificate: %v", err)
		log.Errorln(err)
		return err
	}

	cert = fileChooserCert.GetFilename()
	log.Debugf("selected file to login: %q", cert)
	if cert != "" {
		err = api.API.LoginWithCertsFilename(cert, key, []byte(keypwd))
	} else {
		err = errors.New("no file selected")
	}

	if err != nil {
		v.Dialog(gtk.MESSAGE_ERROR, "An error occured: %v", err)
	}
	return err
}

func (v *Login) onRegisterClicked() (err error) {
	var (
		privkey    string
		privkeypwd string
	)

	fileChooserPrivKey, err := v.FindFileChooserButtonWithBuilder(v.builder, "filechooser_privkey")
	if err != nil {
		err = fmt.Errorf("unable to find file chooser private key: %v", err)
		log.Errorln(err)
		return err
	}
	privkey = fileChooserPrivKey.GetFilename()

	entryPrivKeyPwd, err := v.FindEntryWithBuilder(v.builder, "entry_privpwd")
	if err != nil {
		err = fmt.Errorf("unable to find entry private key password: %v", err)
		log.Errorln(err)
		return err
	}
	privkeypwd, err = entryPrivKeyPwd.GetText()
	if err != nil {
		err = fmt.Errorf("unable to get text from entry private key password: %v", err)
		log.Errorln(err)
		return err
	}

	log.Debugf("selected file to register: %q", privkey)
	if privkey != "" {
		err = api.API.RegisterWithKeyPairFilename(privkey, []byte(privkeypwd))
	} else {
		err = errors.New("no file selected")
	}

	if err != nil {
		v.Dialog(gtk.MESSAGE_ERROR, "An error occured: %v", err)
	}
	return nil
}
