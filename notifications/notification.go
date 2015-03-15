package notifications

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/milou/configuration"
	"github.com/fzzy/radix/redis"
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
	var conn *redis.Client
	var channel string
	var message string
	var marshalled []byte
	var timeout time.Duration

	timeout = time.Duration(conf.Redis.Timeout) * time.Second

	conn, err = redis.DialTimeout("tcp", conf.Redis.Address, timeout)

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

	channel = "milou." + n.Owner + "." + n.Repository

	conn.Cmd("publish", channel, message)

	log.WithFields(log.Fields{
		"message": message,
		"channel": channel,
	}).Info("Published notification")
}
