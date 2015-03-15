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
		http.HandleFunc("/", handlers.Travis)
		http.ListenAndServe(conf.Web.Address, nil)
	} else if conf.Mode == "standalone" {
		for _, project := range projects.List {
			project.Deploy()
		}

		http.HandleFunc("/", handlers.Travis)
		http.ListenAndServe(conf.Web.Address, nil)
	} else if conf.Mode == "client" {
		timeout := time.Duration(10) * time.Second

		conn, err := redis.DialTimeout("tcp", "127.0.0.1:6379", timeout)

		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		psc := pubsub.NewSubClient(conn)

		for _, project := range projects.List {
			_ = psc.Subscribe("milou." + project.Owner + "." + project.Repository)
		}

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
