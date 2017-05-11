package gui

import (
	"errors"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/api"
	"github.com/krostar/nebulo-client-desktop/channel"
	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/gui/view"
)

var baseTitle = "Nebulo - "
var MainWindow *gtk.Window

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
		if err = window.Load(onLoginSucceed); err != nil {
			return fmt.Errorf("unable to build login window: %v", err)
		}
	} else { // login succed, open the main view
		if err = onLoginSucceed(); err != nil {
			return err
		}
	}

	// this block forever until main window is closed
	gtk.Main()
	return nil
}

func onLoginSucceed() (err error) {
	MainWindow := view.Main{}
	MainWindow.WindowBaseTitle = baseTitle
	if err = MainWindow.Load(); err != nil {
		return fmt.Errorf("unable to build main window: %v", err)
	}

	channel.Channels, err = api.API.ChannelList()
	if err != nil {
		MainWindow.Dialog(gtk.MESSAGE_ERROR, "unable to fetch channels list: %v", err)
	}
	err = MainWindow.ChannelsRefresh()
	if err != nil {
		return log.ErrorIf(fmt.Errorf("unable to reresh channel on GUI: %v", err))
	}
	return nil
}
