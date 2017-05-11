package api

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/krostar/nebulo-golib/log"
	"github.com/krostar/nebulo-golib/tools/crypto"

	"github.com/krostar/nebulo-client-desktop/channel"
	"github.com/krostar/nebulo-client-desktop/message"
)

type messageInfos struct {
	Message  message.SecureMsg `json:"message"`
	Receiver string            `json:"receiver_pkey"`
}

type messageCreateRequest struct {
	ChannelName string         `json:"channel_name"`
	Messages    []messageInfos `json:"messages"`
}

func (api *Server) MessageCreate(channelName string, plaintext string) (err error) {
	log.Debugln("doing Message Create call")

	ciphertexts := []messageInfos{}

	for _, member := range channel.Channels[channelName].Members {
		log.Debugf("sign %q with %q", plaintext, member.KeyFingerprint)
		pkeyDER, err := base64.StdEncoding.DecodeString(member.PublicKeyDerBase64)
		if err != nil {
			return fmt.Errorf("failed to decode b64 pkey: %v", err)
		}
		pkey, err := x509.ParsePKIXPublicKey(pkeyDER)
		if err != nil {
			return fmt.Errorf("failed to parse DER encoded public key: %v", err)
		}
		rsaPKey, ok := pkey.(*rsa.PublicKey)
		if !ok {
			return errors.New("cant cast private key to rsa private key")
		}

		ciphertext, keys, hmac, err := crypto.Crypt([]byte(plaintext), *rsaPKey)
		if err != nil {
			return fmt.Errorf("unable to encode message with recipient pkey: %v", err)
		}
		log.Debugln(ciphertext, keys, hmac)
		secureMsg := message.SecureMsg{
			Message:   ciphertext,
			Keys:      keys,
			Integrity: hmac,
		}
		ciphertexts = append(ciphertexts, messageInfos{
			Receiver: member.PublicKeyDerBase64,
			Message:  secureMsg,
		})
	}

	requestBody, err := json.Marshal(&messageCreateRequest{
		ChannelName: channelName,
		Messages:    ciphertexts,
	})
	if err != nil {
		return fmt.Errorf("unable to marshal json: %v", err)
	}

	_, err = api.Post(fmt.Sprintf("chan/%s/message", url.QueryEscape(channelName)), http.StatusCreated, CONTENT_TYPE_JSON, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("unable to get response: %v", err)
	}
	return nil
}
