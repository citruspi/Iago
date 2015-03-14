package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/iago/notifications"
	"github.com/citruspi/iago/webhooks/travis"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TravisWebhook(w http.ResponseWriter, r *http.Request) {
	var announcement travis.Announcement

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&announcement)

	if err != nil {
		log.Error("Failed to decode Travis CI announcement")
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
	} else {
		log.Debug("Retrieved Travis CI authorization")
	}

	log.WithFields(log.Fields{
		"owner":      announcement.Payload.Repository.Owner,
		"repository": announcement.Payload.Repository.Name,
		"branch":     announcement.Payload.Branch,
		"commit":     announcement.Payload.Commit,
	}).Info("Validating Travis CI announcement")

	if !announcement.Valid() {
		log.WithFields(log.Fields{
			"owner":      announcement.Payload.Repository.Owner,
			"repository": announcement.Payload.Repository.Name,
			"branch":     announcement.Payload.Branch,
			"commit":     announcement.Payload.Commit,
		}).Error("Travis CI announcement is invalid")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"owner":      announcement.Payload.Repository.Owner,
		"repository": announcement.Payload.Repository.Name,
		"branch":     announcement.Payload.Branch,
		"commit":     announcement.Payload.Commit,
	}).Debug("Travis CI announcement is valid")

	n := notifications.Build(announcement)
	n.Publish()

	w.WriteHeader(http.StatusOK)
}
