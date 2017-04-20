package gui

import (
	"errors"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/api"
	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/gui/view"
)

var baseTitle = "Nebulo - "

// GUI start the main gui window
func GUI() (err error) {
	gtk.Init(nil)

	var (
		cert   = config.Config.Run.TLS.Cert
		key    = config.Config.Run.TLS.Key
		keypwd = config.Config.Run.TLS.KeyPassword
	)

	// if cert is defined, try to login with it
	if _, err = os.Stat(cert); err == nil {
		if _, err = api.API.LoginWithCertsFilename(cert, key, []byte(keypwd)); err != nil {
			err = fmt.Errorf("unable to log in using %q and %q: %v", cert, key, err)
		}
	} else {
		err = errors.New("cert is undefined or missing")
	}

	if err != nil { // login failed, open the login/register view
		log.Warningf("unable to log in, we have to ask for valid credentials: %v", err)
		window := view.Identity{}
		window.WindowBaseTitle = baseTitle
		if err = window.Load(onLoginSucced); err != nil {
			return fmt.Errorf("unable to build login window: %v", err)
		}
	} else { // login succed, open the main view
		if err = onLoginSucced(); err != nil {
			return err
		}
	}

	// this block forever until main window is closed
	gtk.Main()
	return nil
}

func onLoginSucced() error {
	window := view.Main{}
	window.WindowBaseTitle = baseTitle
	if err := window.Load(); err != nil {
		return fmt.Errorf("unable to build main window: %v", err)
	}
	return nil
}
