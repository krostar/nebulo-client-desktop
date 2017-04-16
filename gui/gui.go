// Package gui ...
// to prevent linter exceptions caused by gtk librairy, theses linters are unused in the whole package: gosimple, errcheck, staticcheck, unused
// nolint: gosimple, errcheck, staticcheck, unused
package gui

import (
	"errors"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/api"
	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/gui/view"
)

// GUI start the main gui window
func GUI() (err error) {
	gtk.Init(nil)

	var (
		baseTitle = "Nebulo - "

		cert   = config.Config.Run.TLS.Cert
		key    = config.Config.Run.TLS.Key
		keypwd = config.Config.Run.TLS.KeyPassword
	)

	// if cert is defined, try to login with it
	if cert != "" {
		if err = api.API.LoginWithCertsFilename(cert, key, []byte(keypwd)); err != nil {
			err = fmt.Errorf("unable to log in using %q and %q: %v", cert, key, err)
		}
	} else {
		err = errors.New("cert is undefined")
	}

	if err != nil {
		log.Warningf("unable to log in, we have to ask for valid credentials: %v", err)
		window := view.Login{}
		window.WindowBaseTitle = baseTitle
		if err = window.Load(); err != nil {
			return fmt.Errorf("unable to build login window: %v", err)
		}
	} else {
		window := view.Main{}
		window.WindowBaseTitle = baseTitle
		if err = window.Load(); err != nil {
			return fmt.Errorf("unable to build main window: %v", err)
		}
	}

	// this block forever until main window is closed
	gtk.Main()
	return nil
}
