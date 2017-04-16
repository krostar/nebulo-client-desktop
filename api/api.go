package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/krostar/nebulo-client-desktop/config"
)

// Server store informations to make the communication
// with the API server easier
type Server struct {
	Client    string
	BaseURL   string
	TLSConfig *tls.Config
	HTTP      *http.Client
}

var API *Server

func createTLSConfig(tlsOptions *config.TLSOptions) (config *tls.Config, err error) {
	clientCAFile, err := ioutil.ReadFile(tlsOptions.ClientsCACert)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %v", clientCAFile, err)
	}
	clientCAPool := x509.NewCertPool()
	clientCAPool.AppendCertsFromPEM(clientCAFile)

	config = &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    clientCAPool,
	}

	if tlsOptions.Cert != "" && tlsOptions.Key != "" {
		crt, err := tls.LoadX509KeyPair(tlsOptions.Cert, tlsOptions.Key)
		if err != nil {
			return nil, fmt.Errorf("unable to load tls key pair: %v", err)
		}
		config.Certificates = []tls.Certificate{crt}
	}

	return config, nil
}

// Get create and send a GET request and return the response
func (api *Server) Get(endpoint string) (response *http.Response, err error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", api.BaseURL, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}
	request.Header.Set("User-Agent", api.Client)

	response, err = api.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %v", err)
	}
	return response, nil
}

// Initialize create a new Server{} base url and certificate configuration
func Initialize(version string, baseurl string, tlsOptions *config.TLSOptions) (serverVersion *VersionResponse, err error) {
	tlsConfig, err := createTLSConfig(tlsOptions)
	if err != nil {
		return nil, fmt.Errorf("tls configuration error: %v", err)
	}

	api := &Server{
		Client:    fmt.Sprintf("nebulo-desktop/%s", version),
		BaseURL:   baseurl,
		TLSConfig: tlsConfig,
		HTTP:      &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}},
	}

	serverVersion, err = api.Version()
	if err != nil {
		return nil, fmt.Errorf("unable to communicate with server %q: %v", baseurl, err)
	}
	API = api
	return serverVersion, nil
}
