package notifications

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/pubsub"
)

type Notification struct {
	Repository string `json:"repository"`
	Owner      string `json:"owner"`
	Commit     string `json:"commit"`
	Branch     string `json:"branch"`
}

var (
	conf configuration.Configuration
)

func init() {
	log.SetLevel(log.DebugLevel)

	conf = configuration.Load()
}

func (n Notification) Publish() {
	var err error
	var channel string
	var message string
	var marshalled []byte

	marshalled, err = json.Marshal(n)

	if err != nil {
		log.Fatal(err)
	}

	message = string(marshalled)

	channel = "milou." + n.Owner + "." + n.Repository

	err = pubsub.Publish(channel, message)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to publish notification")

		return
	}

	log.WithFields(log.Fields{
		"message": message,
		"channel": channel,
	}).Info("Published notification")
}
