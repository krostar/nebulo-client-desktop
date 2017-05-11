package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/channel"
)

type channelListRequest struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

func (api *Server) ChannelList() (list map[string]*channel.Channel, err error) {
	log.Debugln("doing Channel Create call")

	response, err := api.Get("chans", http.StatusOK, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}
	defer response.Body.Close() // nolint: errcheck
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response data: %v", err)
	}
	if err = json.Unmarshal(raw, &list); err != nil {
		return nil, fmt.Errorf("unable to parse response data: %v", err)
	}

	return list, nil
}
