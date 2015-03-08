package host

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	conf "github.com/citruspi/iago/configuration"
	"github.com/citruspi/iago/notification"
	"github.com/citruspi/iago/travis"
)

type Host struct {
	Hostname     string    `json:"hostname"`
	Expiration   time.Time `json:"expiration"`
	Protocol     string    `json:"protocol"`
	Port         int       `json:"port"`
	Path         string    `json:"path"`
	Repositories []string  `json:"repositories"`
}

var (
	List []Host
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func (h Host) URL() string {
	var buffer bytes.Buffer

	buffer.WriteString(h.Protocol)
	buffer.WriteString("://")
	buffer.WriteString(h.Hostname)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(h.Port))
	buffer.WriteString(h.Path)

	log.WithFields(log.Fields{
		"url": string(buffer.Bytes()),
	}).Debug("Prepared host URL")

	return string(buffer.Bytes())
}

func (h Host) CheckIn() {
	log.WithFields(log.Fields{
		"host": h.Hostname,
	}).Info("Checking in host")

	log.WithFields(log.Fields{
		"host": h.Hostname,
	}).Debug("Checking if host is already checked in")

	for i, e := range List {
		if e.Hostname == h.Hostname {

			log.WithFields(log.Fields{
				"host": h.Hostname,
			}).Debug("Host is already checked in")

			diff := List
			diff = append(diff[:i], diff[i+1:]...)
			List = diff

			log.WithFields(log.Fields{
				"host": h.Hostname,
			}).Debug("Removed existing host")
		}
	}

	if h.Protocol == "" {
		log.WithFields(log.Fields{
			"host": h.Hostname,
		}).Debug("Host protocol is undefined")

		h.Protocol = "http"
	}

	if h.Port == 0 {

		log.WithFields(log.Fields{
			"host": h.Hostname,
		}).Debug("Host port is undefined")

		if h.Protocol == "http" {
			h.Port = 80
		} else if h.Protocol == "https" {
			h.Port = 443
		}
	}

	if h.Path == "" {
		log.WithFields(log.Fields{
			"host": h.Hostname,
		}).Debug("Host path is undefined")

		h.Path = "/"
	}

	expiration := time.Now().UTC()
	expiration = expiration.Add(time.Duration(conf.Host.TTL) * time.Second)

	h.Expiration = expiration

	List = append(List, h)

	log.WithFields(log.Fields{
		"host":          h.Hostname,
		"protocol":      h.Protocol,
		"port":          h.Port,
		"path":          h.Path,
		"expiration":    h.Expiration,
		"repositorites": h.Repositories,
	}).Info("Finished checking in host")
}

func Cleanup() {
	for {
		log.Debug("Removing expired hosts")

		for i, h := range List {
			if time.Now().UTC().After(h.Expiration) {
				diff := List
				diff = append(diff[:i], diff[i+1:]...)
				List = diff

				log.WithFields(log.Fields{
					"host": h.Hostname,
				}).Debug("Removed expired host")
			}
		}

		log.Debug("Finished removing expired hosts")

		time.Sleep(time.Duration(conf.Host.TTL) * time.Second)
	}
}

func Notify(announcement travis.Announcement) {
	log.WithFields(log.Fields{
		"repository": announcement.Payload.Repository.Name,
		"owner":      announcement.Payload.Repository.Owner,
		"commit":     announcement.Payload.Commit,
		"branch":     announcement.Payload.Branch,
		"message":    announcement.Payload.Message,
	}).Debug("Preparing to notify hosts about deployment")

	n := notification.Build(announcement)

	log.WithFields(log.Fields{
		"repository": n.Repository,
		"owner":      n.Owner,
		"commit":     n.Commit,
		"branch":     n.Branch,
	}).Debug("Built notification")

	if conf.Notification.Sign {
		n = n.Sign(conf.Notification.PrivateKey)
	}

	content, _ := json.Marshal(n)
	body := bytes.NewBuffer(content)

	var buffer bytes.Buffer

	buffer.WriteString(n.Owner)
	buffer.WriteString("/")
	buffer.WriteString(n.Repository)

	identifier := string(buffer.Bytes())

	for _, host := range List {
		for _, repository := range host.Repositories {
			if repository == identifier {
				log.WithFields(log.Fields{
					"repository": n.Repository,
					"owner":      n.Owner,
					"commit":     n.Commit,
					"branch":     n.Branch,
					"host":       host.Hostname,
				}).Debug("Preparing to notify host")

				urlStr := host.URL()

				req, err := http.NewRequest("POST", urlStr, body)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)

				if err != nil {
					log.WithFields(log.Fields{
						"repository": n.Repository,
						"owner":      n.Owner,
						"commit":     n.Commit,
						"branch":     n.Branch,
						"host":       host.Hostname,
						"error":      err,
					}).Error("Failed to notify host about deployment")
				} else {
					defer resp.Body.Close()

					log.WithFields(log.Fields{
						"repository": n.Repository,
						"owner":      n.Owner,
						"commit":     n.Commit,
						"branch":     n.Branch,
						"host":       host.Hostname,
					}).Info("Notified host about deployment")
				}
			}
		}
	}
}
