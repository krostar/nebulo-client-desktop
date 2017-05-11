package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/krostar/nebulo-golib/log"
)

// VersionResponse is the response format wanted from a /version call
type VersionResponse struct {
	Version string `json:"build_version"`
	Time    string `json:"build_time"`
}

// Version return the server versions informations
func (api *Server) Version() (version *VersionResponse, err error) {
	log.Debugln("doing Version call")

	response, err := api.Get("version", http.StatusOK, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}

	defer response.Body.Close() // nolint: errcheck
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response data: %v", err)
	}

	vr := &VersionResponse{}
	if err = json.Unmarshal(raw, vr); err != nil {
		log.Debugf("ERROR VERSION: %v", err)
		return nil, fmt.Errorf("unable to parse response data: %v", err)
	}

	return vr, nil
}
