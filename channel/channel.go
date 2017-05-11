package channel

import (
	"time"

	"github.com/krostar/nebulo-client-desktop/user"
)

type Channel struct {
	Name             string      `json:"name"`
	Created          time.Time   `json:"created"`
	Creator          user.User   `json:"creator"`
	Members          []user.User `json:"members"`
	MembersCanEdit   bool        `json:"members_can_edit"`
	MembersCanInvite bool        `json:"members_can_invite"`
}

var Channels map[string]*Channel
