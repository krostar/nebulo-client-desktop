package user

import (
	"time"

	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/config"
)

// User represent the informations we get from the API
type User struct {
	KeyFingerprint     string    `json:"key_fingerprint"`
	DisplayName        string    `json:"display_name"`
	Signup             time.Time `json:"signup"`
	LoginFirst         time.Time `json:"login_first"`
	LoginLast          time.Time `json:"login_last"`
	PublicKeyDerBase64 string    `json:"public_key_der_b64"`
}

// Logged store the current logged user
var Logged *User

// Login is called when the login api call succeed
func Login(u *User) {
	if Logged != nil {
		log.Warningln("user already Logged: %q, disconnect first", Logged.KeyFingerprint)
		Logout()
	}
	log.Infof("login successful, Logged user: %q", u.KeyFingerprint)
	Logged = u
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
