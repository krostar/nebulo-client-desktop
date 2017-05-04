package user

import (
	"time"

	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/contact"
)

// User represent the informations we get from the API
type User struct {
	KeyFingerprint     string            `json:"key_fingerprint"`
	DisplayName        string            `json:"display_name"`
	Signup             time.Time         `json:"signup"`
	LoginFirst         time.Time         `json:"login_first"`
	LoginLast          time.Time         `json:"login_last"`
	PublicKeyDerBase64 string            `json:"public_key_der_b64"`
	Contacts           []contact.Contact `json:"contacts"`
}

// Logged store the current logged user
var Logged *User

// Login is called when the login api call succeed
func Login(u *User) (loggedUser *User, err error) {
	if Logged != nil {
		log.Warningln("user already Logged: %q, disconnect first", Logged.KeyFingerprint)
		Logout()
	}
	log.Infof("login successful, Logged user: %q", u.KeyFingerprint)
	if config.Config.Run.ContactsFile != "" {
		u.Contacts, err = contact.LoadFromJSONFile(config.Config.Run.ContactsFile)
		if err != nil {
			log.Warningf("unable to load user contacts from %q: %v, user doesn't have contact yet", config.Config.Run.ContactsFile, err)
		}
	}
	Logged = u
	return Logged, nil
}

// Logout is called when user need to be disconnected
func Logout() {
	if Logged != nil {
		log.Infoln("logout user %q", Logged.KeyFingerprint)
		config.Config.Run.TLS.Key = ""
		config.Config.Run.TLS.KeyPassword = ""
		Logged = nil
	}
}
