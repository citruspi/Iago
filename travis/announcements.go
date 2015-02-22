package travis

import (
	"crypto/sha256"
	"encoding/hex"

	conf "github.com/citruspi/Iago/configuration"
)

type Announcement struct {
	Payload       Payload `json:"payload"`
	Authorization string
}

type Payload struct {
	Status     string     `json:"status_message"`
	Commit     string     `json:"commit"`
	Branch     string     `json:"branch"`
	Message    string     `json:"message"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	Name  string `json:"name"`
	Owner string `json:"owner_name"`
}

func calculateAuthorization(owner string, repository string) string {
	hash := sha256.New()
	hash.Write([]byte(owner))
	hash.Write([]byte("/"))
	hash.Write([]byte(repository))
	hash.Write([]byte(conf.Travis.Token))

	return hex.EncodeToString(hash.Sum(nil))
}

func (a Announcement) Valid() bool {
	if (a.Payload.Status != "Passed") && (a.Payload.Status != "Fixed") {
		return false
	}

	if conf.Travis.Authenticate {
		owner := a.Payload.Repository.Owner
		repository := a.Payload.Repository.Name

		if a.Authorization != calculateAuthorization(owner, repository) {
			return false
		}
	}

	return true
}
