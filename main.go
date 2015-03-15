package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/handlers"
	"github.com/citruspi/milou/notifications"
	"github.com/citruspi/milou/projects"
	"github.com/fzzy/radix/extra/pubsub"
	"github.com/fzzy/radix/redis"
)

var (
	conf configuration.Configuration
)

func init() {
	conf = configuration.Load()
}

func main() {
	if conf.Mode == "server" {
		http.HandleFunc("/webhooks/travis/", handlers.Travis)
		http.ListenAndServe(conf.Web.Address, nil)
	} else if conf.Mode == "standalone" {
		projects.DeployAll()

		http.HandleFunc("/", handlers.Travis)
		http.ListenAndServe(conf.Web.Address, nil)
	} else if conf.Mode == "client" {
		timeout := time.Duration(conf.Redis.Timeout) * time.Second

		conn, err := redis.DialTimeout("tcp", conf.Redis.Address, timeout)

		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		psc := pubsub.NewSubClient(conn)

		projects.DeployAll()

		for {
			psr := psc.Receive()

			if !psr.Timeout() {
				var notification notifications.Notification

				err = json.Unmarshal([]byte(psr.Message), &notification)

				if err != nil {
					log.Fatal(err)
				}

				projects.Process(notification)
			} else {
				continue
			}
		}
	}
}
