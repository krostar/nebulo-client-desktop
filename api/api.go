package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/krostar/nebulo-client-desktop/config"
)

// API store informations to make the communication
// with the API server easier
type API struct {
	Client    string
	BaseURL   string
	TLSConfig *tls.Config
	HTTP      *http.Client
}

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
func (api *API) Get(endpoint string) (response *http.Response, err error) {
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

// New create a new API{} base url and certificate configuration
func New(version string, baseurl string, tlsOptions *config.TLSOptions) (api *API, serverVersion *VersionResponse, err error) {
	tlsConfig, err := createTLSConfig(tlsOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("tls configuration error: %v", err)
	}

	api = &API{
		Client:    fmt.Sprintf("nebulo-desktop/%s", version),
		BaseURL:   baseurl,
		TLSConfig: tlsConfig,
		HTTP:      &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}},
	}

	serverVersion, err = api.Version()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to communicate with server %q: %v", baseurl, err)
	}
	return api, serverVersion, nil
}
