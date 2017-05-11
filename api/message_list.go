package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/message"
)

type messageListRequest struct {
	LastRead time.Time `url:"last_read"`
	Limit    int       `url:"limit"`
}

func (api *Server) MessageList(channelName string, lastRead time.Time) (list []*message.Message, err error) {
	log.Debugln("doing Message List call")

	mlr := &messageListRequest{
		LastRead: lastRead,
		Limit:    -50,
	}
	queryParams, err := query.Values(mlr)
	if err != nil {
		return nil, fmt.Errorf("unable to format query params: %v", err)
	}

	response, err := api.Get(fmt.Sprintf("chan/%s/messages", url.QueryEscape(channelName)), http.StatusOK, queryParams)
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}
	defer response.Body.Close() // nolint: errcheck
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response data: %v", err)
	}

	list = []*message.Message{}
	if err = json.Unmarshal(raw, &list); err != nil {
		return nil, fmt.Errorf("unable to parse response data: %v", err)
	}

	return list, nil
}
