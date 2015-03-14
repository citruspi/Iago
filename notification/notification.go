package notification

import (
	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/iago/travis"
)

type Notification struct {
	Repository string `json:"repository"`
	Owner      string `json:"owner"`
	Commit     string `json:"commit"`
	Branch     string `json:"branch"`
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func Build(announcement travis.Announcement) Notification {
	log.WithFields(log.Fields{
		"repository": announcement.Payload.Repository.Name,
		"owner":      announcement.Payload.Repository.Owner,
		"commit":     announcement.Payload.Commit,
		"branch":     announcement.Payload.Branch,
	}).Debug("Building deployment notification")

	notification := Notification{}
	notification.Repository = announcement.Payload.Repository.Name
	notification.Owner = announcement.Payload.Repository.Owner
	notification.Commit = announcement.Payload.Commit
	notification.Branch = announcement.Payload.Branch

	return notification
}
