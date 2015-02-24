package host

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	conf "github.com/citruspi/iago/configuration"
	"github.com/citruspi/iago/notification"
	"github.com/citruspi/iago/travis"
)

type Host struct {
	Hostname   string    `json:"hostname"`
	Expiration time.Time `json:"expiration"`
	Protocol   string    `json:"protocol"`
	Port       int       `json:"port"`
	Path       string    `json:"path"`
}

var (
	List []Host
)

func (h Host) URL() string {
	var buffer bytes.Buffer

	buffer.WriteString(h.Protocol)
	buffer.WriteString("://")
	buffer.WriteString(h.Hostname)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(h.Port))
	buffer.WriteString(h.Path)

	return string(buffer.Bytes())
}

func (h Host) CheckIn() {
	for i, e := range List {
		if e.Hostname == h.Hostname {
			diff := List
			diff = append(diff[:i], diff[i+1:]...)
			List = diff
		}
	}

	if h.Protocol == "" {
		h.Protocol = "http"
	}

	if h.Port == 0 {
		if h.Protocol == "http" {
			h.Port = 80
		} else if h.Protocol == "https" {
			h.Port = 443
		}
	}

	if h.Path == "" {
		h.Path = "/"
	}

	expiration := time.Now().UTC()
	expiration = expiration.Add(time.Duration(conf.Host.TTL) * time.Second)

	h.Expiration = expiration

	List = append(List, h)
}

func Cleanup() {
	for {
		for i, h := range List {
			if time.Now().UTC().After(h.Expiration) {
				diff := List
				diff = append(diff[:i], diff[i+1:]...)
				List = diff
			}
		}
		time.Sleep(time.Duration(conf.Host.TTL) * time.Second)
	}
}

func Notify(announcement travis.Announcement) {
	n := notification.Build(announcement)

	if conf.Notification.Sign {
		n = n.Sign(conf.Notification.PrivateKey)
	}

	content, _ := json.Marshal(n)
	body := bytes.NewBuffer(content)

	for _, host := range List {
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
}
