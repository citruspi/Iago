package handlers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/iago/notification"
	"github.com/citruspi/iago/webhooks/travis"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TravisWebhook(c *gin.Context) {
	var announcement travis.Announcement

	c.Bind(&announcement)

	log.WithFields(log.Fields{
		"owner":      announcement.Payload.Repository.Owner,
		"repository": announcement.Payload.Repository.Name,
		"branch":     announcement.Payload.Branch,
		"commit":     announcement.Payload.Commit,
	}).Info("Received Travis CI announcement")

	if authorization, ok := c.Request.Header["Authorization"]; ok {
		announcement.Authorization = authorization[0]
		log.Debug("Retrieved Travis CI authorization")
	} else {
		log.Error("Failed to retrieve Travis CI authorization")
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

		c.JSON(400, "")
		return
	} else {
		log.WithFields(log.Fields{
			"owner":      announcement.Payload.Repository.Owner,
			"repository": announcement.Payload.Repository.Name,
			"branch":     announcement.Payload.Branch,
			"commit":     announcement.Payload.Commit,
		}).Debug("Travis CI announcement is valid")

		n := notification.Build(announcement)
		n.Publish()

		c.JSON(200, "")
	}
}
