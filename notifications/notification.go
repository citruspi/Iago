package notifications

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
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
