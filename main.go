package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Host struct {
	Hostname   string    `json:"hostname"`
	Expiration time.Time `json:"expiration"`
}

type Status struct {
	Hosts  []Host  `json:"hosts"`
	Uptime float64 `json:"uptime"`
}

type TravisNotification struct {
	Payload TravisPayload `json:"payload"`
}

type TravisPayload struct {
	Status     string           `json:"status_message"`
	Commit     string           `json:"commit"`
	Branch     string           `json:"branch"`
	Message    string           `json:"message"`
	Repository TravisRepository `json:"repository"`
}

type TravisRepository struct {
	Name  string `json:"name"`
	Owner string `json:"owner_name"`
}

var (
	hosts []Host
	TTL   = 30
)

func buildExpiration() time.Time {
	expiration := time.Now().UTC()
	expiration = expiration.Add(time.Duration(TTL) * time.Second)

	return expiration
}

func cleanup() {
	for {
		for i, h := range hosts {
			if time.Now().UTC().After(h.Expiration) {
				diff := hosts
				diff = append(diff[:i], diff[i+1:]...)
				hosts = diff
			}
		}
		time.Sleep(time.Duration(TTL) * time.Second)
	}
}

func main() {
	boot := time.Now()

	go cleanup()

	router := gin.Default()

	router.POST("/checkin/", func(c *gin.Context) {
		var host Host

		c.Bind(&host)

		if host.Hostname != "" {
			exists := false

			for i, h := range hosts {
				if h.Hostname == host.Hostname {
					exists = true
					hosts[i].Expiration = buildExpiration()
				}
			}

			if !exists {
				host.Expiration = buildExpiration()
				hosts = append(hosts, host)
			}

			c.JSON(200, "")
		} else {
			var response struct {
				Message string `json:"message"`
			}
			response.Message = "Invalid hostname."
			c.JSON(400, response)
		}
	})

	router.GET("/status/", func(c *gin.Context) {
		status := Status{}
		status.Hosts = hosts
		status.Uptime = time.Now().Sub(boot).Seconds()

		c.JSON(200, status)
	})

	router.Run(":8080")
}
