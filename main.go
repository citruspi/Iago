package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/citruspi/Iago/host"
	"github.com/citruspi/Iago/travis"
	"github.com/fogcreek/mini"
	"github.com/gin-gonic/gin"
)

type Status struct {
	Hosts  []host.Host `json:"hosts"`
	Uptime float64     `json:"uptime"`
}

var (
	hosts []host.Host
	TTL   int64
	boot  time.Time
)

func buildExpiration() time.Time {
	expiration := time.Now().UTC()
	expiration = expiration.Add(time.Duration(TTL) * time.Second)

	return expiration
}

func buildURL(h host.Host) string {
	var buffer bytes.Buffer

	buffer.WriteString(h.Protocol)
	buffer.WriteString("://")
	buffer.WriteString(h.Hostname)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(h.Port))
	buffer.WriteString(h.Path)

	return string(buffer.Bytes())
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

func announce(c *gin.Context) {
	var notification travis.Notification

	c.Bind(&notification)

	payload, _ := json.Marshal(notification.Payload)
	body := bytes.NewBuffer(payload)

	for _, host := range hosts {
		urlStr := buildURL(host)

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

		if newHost.Protocol == "" {
			newHost.Protocol = "http"
		}

		if newHost.Port == 0 {
			if newHost.Protocol == "http" {
				newHost.Port = 80
			} else if newHost.Protocol == "https" {
				newHost.Port = 443
			}
		}

		if newHost.Path == "" {
			newHost.Path = "/"
		}

		newHost.Expiration = buildExpiration()
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

	config, err := mini.LoadConfiguration("iago.ini")

	if err != nil {
		log.Fatal(err)
	}

	TTL = config.IntegerFromSection("Hosts", "TTL", 30)

	go cleanup()

	router := gin.Default()

	router.POST("/checkin/", checkin)
	router.POST("/announce/", announce)

	statusEnabled := config.BooleanFromSection("Web", "Status", true)

	if statusEnabled {
		router.GET("/status/", status)
	}

	address := config.StringFromSection("Web", "Address", "127.0.0.1:8080")

	router.Run(address)
}
