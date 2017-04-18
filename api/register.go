package api

import (
	"crypto"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"github.com/krostar/nebulo-golib/log"
	"github.com/krostar/nebulo-golib/tools/cert"
)

// Register send a certificate signing request and store the signed certificate
func (api *Server) Register(certificateRequest *x509.CertificateRequest) (err error) {
	log.Debugln("doing Register call")
	//sign certificate
	signedCertificate, err := cert.ParseCSR(certificateRequest)
	if err != nil {
		return fmt.Errorf("Unable to sign certificate",err)
	}
	
	//encode certificate
	encodedPEMCertByte, err := cert.EncodeCertificatePEM(signedCertificate)
	if err != nil {
		return fmt.Errorf("Unable to encode certificate",err)
	}

	//send back to client save to file
	err = ioutil.WriteFile("LoginCertificate", encodedPEMCertByte, 0644)
	if err != nil{
		return fmt.Errorf("Unable to write to file", err)
	}
	file, err := os.Create("LoginCertificate")
	defer file.Close()

	return errors.New("unhandled")
}

// RegisterWithKeyPairFilename do the same thing as Register but with key path and password
func (api *Server) RegisterWithKeyPairFilename(privateKeyFilepath string, privateKeyPassword []byte, emailAddress string, commonName string) (err error) {
	key, err := cert.ParsePrivateKeyPEMFromFile(privateKeyFilepath, privateKeyPassword)
	if err != nil {
		return fmt.Errorf("unable to get certificates from file: %v", err)
	}
	
	certificateRequest, err := api.CreateCertificateRequest(key, emailAddress, commonName)
	if err != nil {
		return fmt.Errorf("unable to get certificate request: %v", err)
	}
	 
	return api.Register(certificateRequest)
}

//Create Certificate Request
func (api *Server) CreateCertificateRequest(key crypto.Signer, emailAddress string, commonName string) (csr *x509.CertificateRequest, err error) {
	var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
	 
	subj := pkix.Name{
		CommonName:         commonName,
		Country:            []string{"US"},
		Province:           []string{"Some-State"},
		Locality:           []string{"MyCity"},
		Organization:       []string{"Company Ltd"},
		OrganizationalUnit: []string{"IT"},
	}
		
	rawSubj := subj.ToRDNSequence()
	rawSubj = append(rawSubj, []pkix.AttributeTypeAndValue{
		{Type: oidEmailAddress, Value: emailAddress},
	})

	asn1Subj, _ := asn1.Marshal(rawSubj)
	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		EmailAddresses:     []string{emailAddress},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	
	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, key)
	return x509.ParseCertificateRequest(csrBytes)
	

	
}