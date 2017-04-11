package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/krostar/nebulo-golib/log"
)

// VersionResponse is the response format wanted from a /version call
type VersionResponse struct {
	Version string `json:"build_version"`
	Time    string `json:"build_time"`
}

// Version return the server versions informations
func (api *API) Version() (version *VersionResponse, err error) {
	log.Debugln("doing Version call")
	response, err := api.Get("version")
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			panic(err)
		}
	}()
	log.Debugln("response status: %s", response.Status)

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response data: %v", err)
	}

	vr := &VersionResponse{}
	if err = json.Unmarshal(raw, vr); err != nil {
		return nil, fmt.Errorf("unable to parse response data: %v", err)
	}

	return vr, nil
}
