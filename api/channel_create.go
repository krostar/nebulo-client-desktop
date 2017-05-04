package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/channel"
)

// ChannelCreateResponse is the response format wanted from a /channel call
type ChannelCreateResponse channel.Channel

type channelCreateRequest struct {
	Name             string   `json:"name"`
	MembersPublicKey []string `json:"members_public_key"`
}

// ChannelCreate return the wanted channel profile informations
func (api *Server) ChannelCreate(name string, membersPublicKey []string) (c *channel.Channel, err error) {
	log.Debugln("doing Channel Create call")

	requestBody, err := json.Marshal(&channelCreateRequest{
		Name:             name,
		MembersPublicKey: membersPublicKey,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json: %v", err)
	}

	response, err := api.Post("chan", http.StatusOK, CONTENT_TYPE_JSON, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}
	defer response.Body.Close() // nolint: errcheck
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response data: %v", err)
	}

	crr := &ChannelCreateResponse{}
	if err = json.Unmarshal(raw, crr); err != nil {
		return nil, fmt.Errorf("unable to parse response data: %v", err)
	}

	return (*channel.Channel)(crr), nil
}
