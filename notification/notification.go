package notification

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"io/ioutil"
	"log"

	"github.com/citruspi/Iago/travis"
)

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

func (notification Notification) Sign(privateKeyPath string) Notification {
	privateKeyRaw, err := ioutil.ReadFile(privateKeyPath)

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

func Build(announcement travis.Announcement) Notification {
	notification := Notification{}
	notification.Repository = announcement.Payload.Repository.Name
	notification.Owner = announcement.Payload.Repository.Owner
	notification.Commit = announcement.Payload.Commit
	notification.Branch = announcement.Payload.Branch

	return notification
}
