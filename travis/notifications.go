package travis

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/citruspi/Iago/configuration"
)

type Notification struct {
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
	hash.Write([]byte(configuration.Conf.Travis.Token))

	return hex.EncodeToString(hash.Sum(nil))
}

func (n Notification) Valid() bool {
	if (n.Payload.Status != "Passed") && (n.Payload.Status != "Fixed") {
		return false
	}

	if configuration.Conf.Travis.Authenticate {
		owner := n.Payload.Repository.Owner
		repository := n.Payload.Repository.Name

		if n.Authorization != calculateAuthorization(owner, repository) {
			return false
		}
	}

	return true
}
