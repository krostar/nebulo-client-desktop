package api

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/krostar/nebulo-golib/log"
	"github.com/krostar/nebulo-golib/tools/cert"

	"github.com/krostar/nebulo-client-desktop/config"
	"github.com/krostar/nebulo-client-desktop/user"
)

// Register send a certificate signing request and store the signed certificate
func (api *Server) Register(key crypto.PrivateKey) (newUser *user.User, err error) {
	log.Debugln("doing Register call")

	if user.Logged != nil {
		return nil, errors.New("user already logged, no need to register")
	}

	// first, create a certificate signing request from key
	csr, err := api.createCSR(key)
	if err != nil {
		return nil, err
	}

	// send that csr to the server and pray he accept to sign it
	response, err := api.Post("user", http.StatusCreated, "application/x-pem-file", bytes.NewReader(csr))
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}

	// the response is supposed to contain the signed certificate
	defer response.Body.Close() // nolint: errcheck
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response data: %v", err)
	}

	// save it
	if err = ioutil.WriteFile(config.Config.Run.TLS.Cert, raw, 0400); err != nil {
		return nil, fmt.Errorf("unable to write identity cert file %q: %v", config.Config.Run.TLS.Cert, err)
	}
	if err = config.SaveFile(); err != nil {
		return nil, fmt.Errorf("unable to save configuration file: %v", err)
	}

	// try to log with this certificate
	return api.Login()
}

func (api *Server) createCSR(key crypto.PrivateKey) (_ []byte, err error) {
	subj := pkix.Name{
		CommonName:         "nebulo-client",
		Country:            []string{"-"},
		Province:           []string{"-"},
		Locality:           []string{"-"},
		Organization:       []string{"Nebulo"},
		OrganizationalUnit: []string{"Nebulo Clients CA"},
	}

	asn1Subj, err := asn1.Marshal(subj.ToRDNSequence())
	if err != nil {
		return nil, fmt.Errorf("unable to marshal asn1: %v", err)
	}

	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		SignatureAlgorithm: x509.SHA512WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, key)
	if err != nil {
		return nil, fmt.Errorf("unable to create csr: %v", err)
	}

	var b bytes.Buffer
	bb := bufio.NewWriter(&b)
	if err = pem.Encode(bb, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}); err != nil {
		return nil, fmt.Errorf("unable to encode pem: %v", err)
	}
	if err = bb.Flush(); err != nil {
		return nil, fmt.Errorf("unable to flush buffer: %v", err)
	}
	return b.Bytes(), nil
}

// RegisterWithKeyPairFilename do the same thing as Register but with key path and password
func (api *Server) RegisterWithKeyPairFilename(privateKeyFilepath string, privateKeyPassword []byte) (_ *user.User, err error) {
	key, err := cert.ParsePrivateKeyPEMFromFile(privateKeyFilepath, privateKeyPassword)
	if err != nil {
		return nil, fmt.Errorf("unable to get key from file: %v", err)
	}

	config.Config.Run.TLS.Key = privateKeyFilepath
	config.Config.Run.TLS.KeyPassword = string(privateKeyPassword)
	if err = changeTLSOptions(API, &config.Config.Run.TLS); err != nil {
		return nil, fmt.Errorf("unable to change tls options to register: %v", err)
	}

	return api.Register(key)
}
