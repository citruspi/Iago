package notification

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/iago/travis"
	"github.com/fzzy/radix/redis"
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

func (n Notification) Publish() {
	var err error
	var conn *redis.Client
	var channel string
	var message string
	var marshalled []byte
	var timeout time.Duration

	timeout = time.Duration(10) * time.Second

	conn, err = redis.DialTimeout("tcp", "127.0.0.1:6379", timeout)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Established connection to redis")
	}

	defer conn.Close()

	marshalled, err = json.Marshal(n)

	if err != nil {
		log.Fatal(err)
	}

	message = string(marshalled)

	channel = "iago." + n.Owner + "." + n.Repository

	conn.Cmd("publish", channel, message)

	log.WithFields(log.Fields{
		"message": message,
		"channel": channel,
	}).Info("Published notification")
}
