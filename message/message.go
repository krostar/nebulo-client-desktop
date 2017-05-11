package message

import (
	"time"

	"github.com/krostar/nebulo-client-desktop/channel"
	"github.com/krostar/nebulo-client-desktop/user"
)

type SecureMsg struct {
	Message   []byte `json:"message"`
	Keys      []byte `json:"keys"`
	Integrity []byte `json:"integrity"`
}

type Message struct {
	Ciphertext []byte          `json:"message"`
	Keys       []byte          `json:"keys"`
	Integrity  []byte          `json:"integrity"`
	Plaintext  string          `json:"plaintext"`
	Channel    channel.Channel `json:"channel"`
	Sender     user.User       `json:"sender"`
	Posted     time.Time       `json:"posted"`
}
