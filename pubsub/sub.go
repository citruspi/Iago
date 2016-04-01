package pubsub

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fzzy/radix/extra/pubsub"
	"github.com/fzzy/radix/redis"
)

func Subscribe(channels []string) (*pubsub.SubClient, error) {
	var err error
	var conn *redis.Client
	var timeout time.Duration
	var client *pubsub.SubClient

	timeout = time.Duration(conf.Redis.Timeout) * time.Second

	conn, err = redis.DialTimeout("tcp", conf.Redis.Address, timeout)

	if err != nil {
		return client, err
	} else {
		log.Info("Established connection to redis")
	}

	client = pubsub.NewSubClient(conn)

	for _, channel := range channels {
		r := client.Subscribe(channel)

		if r.Err != nil {
			return client, r.Err
		}
	}

	return client, nil
}
