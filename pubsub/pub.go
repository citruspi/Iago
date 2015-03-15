package pubsub

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/milou/configuration"
	"github.com/fzzy/radix/redis"
)

var (
	conf configuration.Configuration
)

func init() {
	conf = configuration.Load()
}

func Publish(channel string, message string) error {
	var err error
	var conn *redis.Client
	var timeout time.Duration

	timeout = time.Duration(conf.Redis.Timeout) * time.Second

	conn, err = redis.DialTimeout("tcp", conf.Redis.Address, timeout)

	if err != nil {
		return err
	} else {
		log.Info("Established connection to redis")
	}

	defer conn.Close()

	conn.Cmd("publish", channel, message)

	return nil
}
