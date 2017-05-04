package contact

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/krostar/nebulo-golib/log"
)

// Contact store a user contact
type Contact struct {
	Name         string `json:"name"`
	PublicKeyB64 string `json:"public_key_b64"`
}

func LoadFromJSONFile(filepath string) (contacts []Contact, err error) {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %q: %v", filepath, err)
	}

	contacts = []Contact{}

	if err = json.Unmarshal(raw, &contacts); err != nil {
		return nil, fmt.Errorf("unable to parse json file: %v", err)
	}

	return contacts, nil
}

func AddToFile(filepath string, name string, publicKeyB64 string) (contacts []Contact, err error) {
	contacts, err = LoadFromJSONFile(filepath)
	if err != nil {
		log.Warningf("unable to load user contacts from %q: %v, save will replace file content", filepath, err)
		contacts = []Contact{}
	}

	newContact := Contact{
		Name:         name,
		PublicKeyB64: publicKeyB64,
	}
	contacts = append(contacts, newContact)

	contactsJSON, err := json.MarshalIndent(contacts, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("unable to create json: %v", err)
	}
	if err := ioutil.WriteFile(filepath, contactsJSON, 0600); err != nil {
		return nil, fmt.Errorf("unable to write configuration file %q: %v", filepath, err)
	}
	return contacts, nil
}
