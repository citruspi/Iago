package travis

import (
	"crypto/sha256"
	"encoding/hex"

	log "github.com/Sirupsen/logrus"
	conf "github.com/citruspi/iago/configuration"
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

func init() {
	log.SetLevel(log.DebugLevel)
}

func calculateAuthorization(owner string, repository string) string {
	log.WithFields(log.Fields{
		"owner":      owner,
		"repository": repository,
		"token":      conf.Travis.Token,
	}).Debug("Calculating Travis CI authorization")

	hash := sha256.New()
	hash.Write([]byte(owner))
	hash.Write([]byte("/"))
	hash.Write([]byte(repository))
	hash.Write([]byte(conf.Travis.Token))

	authorization := hex.EncodeToString(hash.Sum(nil))

	log.WithFields(log.Fields{
		"owner":         owner,
		"repository":    repository,
		"token":         conf.Travis.Token,
		"authorization": authorization,
	}).Debug("Calculated Travis CI authorization")

	return authorization
}

func (a Announcement) Valid() bool {
	if (a.Payload.Status != "Passed") && (a.Payload.Status != "Fixed") {
		log.WithFields(log.Fields{
			"status":     a.Payload.Status,
			"branch":     a.Payload.Branch,
			"commit":     a.Payload.Commit,
			"owner":      a.Payload.Repository.Owner,
			"repository": a.Payload.Repository.Name,
		}).Error("Announcement with unacceptable status received")

		return false
	}

	if conf.Travis.Authenticate {
		owner := a.Payload.Repository.Owner
		repository := a.Payload.Repository.Name

		if a.Authorization != calculateAuthorization(owner, repository) {
			log.WithFields(log.Fields{
				"received":   a.Authorization,
				"calculated": calculateAuthorization(owner, repository),
			}).Error("Incorrect authorization received")

			return false
		}
	}

	return true
}
