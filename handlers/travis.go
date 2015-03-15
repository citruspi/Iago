package handlers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/projects"
	"github.com/citruspi/milou/webhooks/travis"
)

var (
	conf configuration.Configuration
)

func init() {
	log.SetLevel(log.DebugLevel)

	conf = configuration.Load()
}

func Travis(w http.ResponseWriter, r *http.Request) {
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

	notification := announcement.ToNotification()

	if conf.Mode == "server" {
		notification.Publish()
	} else if conf.Mode == "standalone" {
		projects.Process(notification)
	}

	w.WriteHeader(http.StatusOK)
}
