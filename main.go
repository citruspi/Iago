package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/handlers"
	"github.com/citruspi/milou/notifications"
	"github.com/citruspi/milou/projects"
	"github.com/citruspi/milou/pubsub"
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
		projects.DeployAll()

		var channels []string

		for _, project := range projects.Projects {
			channels = append(channels, "milou."+project.Owner+"."+project.Repository)
		}

		client, err := pubsub.Subscribe(channels)

		if err != nil {
			log.Fatal(err)
		}

		for {
			reply := client.Receive()

			if !reply.Timeout() {
				var notification notifications.Notification

				err = json.Unmarshal([]byte(reply.Message), &notification)

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
