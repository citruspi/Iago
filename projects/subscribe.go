package projects

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/milou/notifications"
	"github.com/citruspi/milou/pubsub"
)

func Subscribe() {
	var channel string
	var channels []string

	for _, project := range list {
		channel = "milou." + project.Owner + "." + project.Repository
		channels = append(channels, channel)
	}

	client, err := pubsub.Subscribe(channels)

	if err != nil {
		log.Fatal(err)
	}

	for {
		reply := client.Receive()

		if reply.Timeout() {
			continue
		} else if reply.Err != nil {
			log.Error(reply.Err)
		} else {
			var notification notifications.Notification

			err = json.Unmarshal([]byte(reply.Message), &notification)

			if err != nil {
				log.Fatal(err)
			}

			Process(notification)
		}
	}
}
