package handlers

import (
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

	announcement = travis.ProcessRequest(r)

	if !announcement.Valid {
		log.WithFields(log.Fields{
			"owner":      announcement.Payload.Repository.Owner,
			"repository": announcement.Payload.Repository.Name,
			"branch":     announcement.Payload.Branch,
			"commit":     announcement.Payload.Commit,
		}).Error("Travis CI announcement is invalid")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !announcement.Authentic {
		log.WithFields(log.Fields{
			"owner":      announcement.Payload.Repository.Owner,
			"repository": announcement.Payload.Repository.Name,
			"branch":     announcement.Payload.Branch,
			"commit":     announcement.Payload.Commit,
		}).Error("Travis CI announcement is invalid")

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.WithFields(log.Fields{
		"owner":      announcement.Payload.Repository.Owner,
		"repository": announcement.Payload.Repository.Name,
		"branch":     announcement.Payload.Branch,
		"commit":     announcement.Payload.Commit,
	}).Debug("Travis CI announcement is valid and authentic")

	n := notifications.Build(announcement)
	n.Publish()

	w.WriteHeader(http.StatusOK)
}
