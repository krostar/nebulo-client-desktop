package api

import (
	"errors"
	"fmt"

	"github.com/krostar/nebulo-golib/log"
	"github.com/krostar/nebulo-golib/tools/cert"

	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/user"
)

// Login log a user based on his nebulo signed certificate
func (api *Server) Login() (loggedUser *user.User, err error) {
	log.Debugln("doing Login call")

	// there is no login call, just check if the current configuration allow a required-auth call
	loggedUser, err = api.UserProfile()
	if err != nil { // no ? delete current configuration
		config.Config.Run.TLS.Key = ""
		config.Config.Run.TLS.KeyPassword = ""
		err = fmt.Errorf("unable to login: %v", err)
	}
	if errSave := config.SaveFile(); errSave != nil {
		if err == nil {
			err = errors.New("")
		}
		err = fmt.Errorf("%s and unable to save configuration file: %v", err, errSave)
	}

	if err != nil {
		return nil, err
	}

	// login succeed
	user.Login(loggedUser)

	return loggedUser, nil
}

// LoginWithCertsFilename do the Login call but with the cert and key path
func (api *Server) LoginWithCertsFilename(certFilepath string, keyFilePath string, keyPassword []byte) (_ *user.User, err error) {
	_, _, err = cert.CertAndKeyFromFiles(certFilepath, keyFilePath, keyPassword)
	if err != nil {
		return nil, fmt.Errorf("unable to get certificate from file: %v", err)
	}

	config.Config.Run.TLS.Cert = certFilepath
	config.Config.Run.TLS.Key = keyFilePath
	config.Config.Run.TLS.KeyPassword = string(keyPassword)
	if err = changeTLSOptions(API, &config.Config.Run.TLS); err != nil {
		return nil, fmt.Errorf("unable to change tls options to login: %v", err)
	}

	return api.Login()
}
