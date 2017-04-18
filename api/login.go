package api

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/krostar/nebulo-golib/log"
	"github.com/krostar/nebulo-golib/tools/cert"
)

// Login log a user based on his nebulo signed certificate
func (api *Server) Login(cert *x509.Certificate, key crypto.Signer) (err error) {
	log.Debugln("doing Login call")

	//verify if certificate has been revoked
	revoked, err := cert.VerifyCertificate(cert)

	if err != nil {
		fmt.Errorf("Unable to verfiy certificate.", err)
	}

	if (revoked)  {
		return fmt.Errorf("Your certificate has been revoked.")
	} else {
		fmt.Println("You can log in")
	}

	return errors.New("unhandled")
}

// LoginWithCertsFilename do the Login call but with the cert and key path
func (api *Server) LoginWithCertsFilename(certFilepath string, keyFilepath string, keyPassword []byte) (err error) {
	cert, key, err := cert.CertAndKeyFromFiles(certFilepath, keyFilepath, keyPassword)
	if err != nil {
		return fmt.Errorf("unable to get certificates from file: %v", err)
	}
	return api.Login(cert, key)
}
