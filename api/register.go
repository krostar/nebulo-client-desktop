package api

import (
	"crypto"
	"errors"
	"fmt"

	"github.com/krostar/nebulo-golib/log"
	"github.com/krostar/nebulo-golib/tools/cert"
)

// Register send a certificate signing request and store the signed certificate
func (api *Server) Register(key crypto.Signer) (err error) {
	log.Debugln("doing Register call")
	return errors.New("unhandled")
}

// RegisterWithKeyPairFilename do the same thing as Register but with key path and password
func (api *Server) RegisterWithKeyPairFilename(privateKeyFilepath string, privateKeyPassword []byte) (err error) {
	key, err := cert.ParsePrivateKeyPEMFromFile(privateKeyFilepath, privateKeyPassword)
	if err != nil {
		return fmt.Errorf("unable to get certificates from file: %v", err)
	}
	return api.Register(key)
}
