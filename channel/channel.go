package channel

type Channel struct {
	Name             string   `json:"name"`
	MembersPublicKey []string `json:"members_public_key"`
}
