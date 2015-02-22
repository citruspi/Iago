package host

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	conf "github.com/citruspi/Iago/configuration"
	"github.com/citruspi/Iago/travis"
)

type Host struct {
	Hostname   string    `json:"hostname"`
	Expiration time.Time `json:"expiration"`
	Protocol   string    `json:"protocol"`
	Port       int       `json:"port"`
	Path       string    `json:"path"`
}

type Notification struct {
	Repository string    `json:"repository"`
	Owner      string    `json:"owner"`
	Commit     string    `json:"commit"`
	Branch     string    `json:"branch"`
	Signature  Signature `json:"signature,omitempty"`
}

type Signature struct {
	S string `json:"s,omitempty"`
	R string `json:"r,omitempty"`
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

func (notification Notification) Sign() Notification {
	privateKeyRaw, err := ioutil.ReadFile(conf.Notification.PrivateKey)

	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := x509.ParseECPrivateKey(privateKeyRaw)

	if err != nil {
		log.Fatal(err)
	}

	var message bytes.Buffer
	message.WriteString(notification.Owner)
	message.WriteString("/")
	message.WriteString(notification.Repository)
	message.WriteString("@")
	message.WriteString(notification.Branch)
	message.WriteString("#")
	message.WriteString(notification.Commit)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, message.Bytes())

	if err != nil {
		log.Fatal(err)
	}

	notification.Signature.R = r.String()
	notification.Signature.S = s.String()

	return notification
}

func Notify(announcement travis.Announcement) {
	notification := Notification{}
	notification.Repository = announcement.Payload.Repository.Name
	notification.Owner = announcement.Payload.Repository.Owner
	notification.Commit = announcement.Payload.Commit
	notification.Branch = announcement.Payload.Branch

	if conf.Notification.Sign {
		notification = notification.Sign()
	}

	content, _ := json.Marshal(notification)
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
