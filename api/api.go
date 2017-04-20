package api

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/krostar/nebulo-golib/log"
	ghttperror "github.com/krostar/nebulo-golib/router/httperror"
	"github.com/krostar/nebulo-golib/tools/cert"

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

// API is the current configuration to contact the api server
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

	// if we don't have client cert or private key, nothing to do, call who need auth will failed
	if tlsOptions.Cert != "" && tlsOptions.Key != "" {
		crt, err := cert.TLSCertificateFromFiles(tlsOptions.Cert, tlsOptions.Key, []byte(tlsOptions.KeyPassword))
		if err != nil {
			log.Warningf("unable to load tls key pair: %v", err)
		} else {
			config.Certificates = []tls.Certificate{*crt}
		}
	}

	return config, nil
}

// Request add things every requests need, do the request, check the status code and return the response
func (api *Server) Request(request *http.Request, expectedStatus int) (response *http.Response, err error) {
	request.Header.Set("User-Agent", api.Client)

	response, err = api.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %v", err)
	}

	// if response status doesnt match expected status, read the response, parse the error and return it
	if response.StatusCode != expectedStatus {
		err = fmt.Errorf("bad status code: expected %d receive %d", expectedStatus, response.StatusCode)
		defer response.Body.Close() // nolint: errcheck
		raw, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			err = fmt.Errorf("%v; unable to read response data: %v", err, errRead)
			return nil, err
		}
		er := &ghttperror.HTTPErrors{}
		if errJSON := json.Unmarshal(raw, er); errJSON != nil {
			return nil, fmt.Errorf("%v; unable to parse response data: %v", err, errJSON)
		}
		return nil, fmt.Errorf("%v; %s", err, er.Error())
	}
	return response, nil
}

// Get create and send a GET request and return the response
func (api *Server) Get(endpoint string, expectedStatus int) (response *http.Response, err error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", api.BaseURL, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}
	return api.Request(request, expectedStatus)
}

// Post create and send a POST request and return the response
func (api *Server) Post(endpoint string, expectedStatus int, contentType string, body io.Reader) (response *http.Response, err error) {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", api.BaseURL, endpoint), body)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}
	request.Header.Set("Content-Type", contentType)
	return api.Request(request, expectedStatus)
}

// Initialize create a new Server{} base url and certificate configuration
func Initialize(version string, baseurl string, tlsOptions *config.TLSOptions) (serverVersion *VersionResponse, err error) {
	api := &Server{
		Client:  fmt.Sprintf("nebulo-desktop/%s", version),
		BaseURL: baseurl,
	}
	if err = changeTLSOptions(api, tlsOptions); err != nil {
		return nil, err
	}

	serverVersion, err = api.Version()
	if err != nil {
		return nil, fmt.Errorf("unable to communicate with server %q: %v", baseurl, err)
	}

	API = api
	return serverVersion, nil
}

func changeTLSOptions(api *Server, tlsOptions *config.TLSOptions) (err error) {
	tlsConfig, err := createTLSConfig(tlsOptions)
	if err != nil {
		return fmt.Errorf("tls configuration error: %v", err)
	}

	api.TLSConfig = tlsConfig
	api.HTTP = &http.Client{
		Timeout:   30 * time.Second,
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}
	return nil
}
