package travis

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	conf "github.com/citruspi/iago/configuration"
	"github.com/citruspi/iago/notifications"
)

type Announcement struct {
	Payload       Payload `json:"payload"`
	Authorization string
	Valid         bool
	Authentic     bool
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

func (announcement Announcement) ToNotification() notifications.Notification {
	log.WithFields(log.Fields{
		"repository": announcement.Payload.Repository.Name,
		"owner":      announcement.Payload.Repository.Owner,
		"commit":     announcement.Payload.Commit,
		"branch":     announcement.Payload.Branch,
	}).Debug("Building deployment notification")

	notification := notifications.Notification{}
	notification.Repository = announcement.Payload.Repository.Name
	notification.Owner = announcement.Payload.Repository.Owner
	notification.Commit = announcement.Payload.Commit
	notification.Branch = announcement.Payload.Branch

	return notification
}

func ProcessRequest(r *http.Request) Announcement {
	var announcement Announcement

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&announcement)

	if err != nil {
		log.Error("Failed to decode Travis CI announcement")
		return announcement
	}

	log.WithFields(log.Fields{
		"owner":      announcement.Payload.Repository.Owner,
		"repository": announcement.Payload.Repository.Name,
		"branch":     announcement.Payload.Branch,
		"commit":     announcement.Payload.Commit,
	}).Info("Received Travis CI announcement")

	announcement.Authorization = r.Header.Get("Authorization")

	if announcement.Authorization == "" {
		log.Error("Failed to retrieve Travis CI authorization")
		return announcement
	}

	log.Debug("Retrieved Travis CI authorization")

	log.WithFields(log.Fields{
		"owner":      announcement.Payload.Repository.Owner,
		"repository": announcement.Payload.Repository.Name,
		"branch":     announcement.Payload.Branch,
		"commit":     announcement.Payload.Commit,
	}).Info("Validating Travis CI announcement")

	if (announcement.Payload.Status != "Passed") && (announcement.Payload.Status != "Fixed") {
		log.WithFields(log.Fields{
			"status":     announcement.Payload.Status,
			"branch":     announcement.Payload.Branch,
			"commit":     announcement.Payload.Commit,
			"owner":      announcement.Payload.Repository.Owner,
			"repository": announcement.Payload.Repository.Name,
		}).Error("Announcement with unacceptable status received")

		return announcement
	}

	if (announcement.Payload.Branch == "") ||
		(announcement.Payload.Commit == "") ||
		(announcement.Payload.Repository.Owner == "") ||
		(announcement.Payload.Repository.Name == "") {

		log.WithFields(log.Fields{
			"status":     announcement.Payload.Status,
			"branch":     announcement.Payload.Branch,
			"commit":     announcement.Payload.Commit,
			"owner":      announcement.Payload.Repository.Owner,
			"repository": announcement.Payload.Repository.Name,
		}).Error("Incomplete announcement received")

		return announcement
	}

	announcement.Valid = true

	if conf.Travis.Authenticate {
		owner := announcement.Payload.Repository.Owner
		repository := announcement.Payload.Repository.Name

		if announcement.Authorization != calculateAuthorization(owner, repository) {
			log.WithFields(log.Fields{
				"received":   announcement.Authorization,
				"calculated": calculateAuthorization(owner, repository),
			}).Error("Incorrect authorization received")
		} else {
			announcement.Authentic = true
		}
	}

	return announcement
}
