package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/citruspi/Iago/configuration"
	"github.com/citruspi/Iago/host"
	"github.com/citruspi/Iago/travis"
	"github.com/gin-gonic/gin"
)

type Status struct {
	Hosts  []host.Host `json:"hosts"`
	Uptime float64     `json:"uptime"`
}

var (
	hosts []host.Host
	boot  time.Time
	conf  configuration.Configuration
)

func cleanup() {
	for {
		for i, h := range hosts {
			if time.Now().UTC().After(h.Expiration) {
				diff := hosts
				diff = append(diff[:i], diff[i+1:]...)
				hosts = diff
			}
		}
		time.Sleep(time.Duration(conf.Host.TTL) * time.Second)
	}
}

func announce(c *gin.Context) {
	var notification travis.Notification

	c.Bind(&notification)

	payload, _ := json.Marshal(notification.Payload)
	body := bytes.NewBuffer(payload)

	for _, host := range hosts {
		urlStr := host.URL()

		req, err := http.NewRequest("POST", urlStr, body)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
	}

	c.JSON(200, "")
}

func checkin(c *gin.Context) {
	var newHost host.Host

	c.Bind(&newHost)

	if newHost.Hostname != "" {
		for i, h := range hosts {
			if h.Hostname == newHost.Hostname {
				diff := hosts
				diff = append(diff[:i], diff[i+1:]...)
				hosts = diff
			}
		}

		newHost = newHost.Process()

		hosts = append(hosts, newHost)

		c.JSON(200, "")
	} else {
		var response struct {
			Message string `json:"message"`
		}
		response.Message = "Invalid hostname."
		c.JSON(400, response)
	}
}

func status(c *gin.Context) {
	status := Status{}
	status.Hosts = hosts
	status.Uptime = time.Now().Sub(boot).Seconds()

	c.JSON(200, status)
}

func main() {
	boot = time.Now()

	configuration.Process()
	conf = configuration.Conf

	go cleanup()

	router := gin.Default()

	router.POST("/checkin/", checkin)
	router.POST("/announce/", announce)

	if conf.Web.Status {
		router.GET("/status/", status)
	}

	router.Run(conf.Web.Address)
}
