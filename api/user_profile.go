package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/krostar/nebulo-golib/log"

	"github.com/krostar/nebulo-client-desktop/user"
)

// UserProfileResponse is the response format wanted from a /user call
type UserProfileResponse user.User

// UserProfile return the user profile informations
func (api *Server) UserProfile() (u *user.User, err error) {
	log.Debugln("doing User Profile call")

	response, err := api.Get("user", http.StatusOK, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}
	defer response.Body.Close() // nolint: errcheck
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response data: %v", err)
	}

	upr := &UserProfileResponse{}
	if err = json.Unmarshal(raw, upr); err != nil {
		return nil, fmt.Errorf("unable to parse response data: %v", err)
	}

	return (*user.User)(upr), nil
}
