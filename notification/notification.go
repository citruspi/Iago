package notification

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"io/ioutil"
	"math/big"

	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/iago/travis"
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

func init() {
	log.SetLevel(log.DebugLevel)
}

func (notification Notification) Sign(privateKeyPath string) Notification {
	log.WithFields(log.Fields{
		"owner":      notification.Owner,
		"repository": notification.Repository,
		"branch":     notification.Branch,
		"commit":     notification.Commit,
	}).Debug("Signing notification")

	privateKeyRaw, err := ioutil.ReadFile(privateKeyPath)

	if err != nil {
		log.WithFields(log.Fields{
			"path":  privateKeyPath,
			"error": err,
		}).Fatal("Failed to read private key")
	} else {
		log.WithFields(log.Fields{
			"path": privateKeyPath,
		}).Debug("Read private key")
	}

	privateKey, err := x509.ParseECPrivateKey(privateKeyRaw)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to parse private key")
	} else {
		log.Debug("Parsed private key")
	}

	var message bytes.Buffer
	message.WriteString(notification.Owner)
	message.WriteString("/")
	message.WriteString(notification.Repository)
	message.WriteString("@")
	message.WriteString(notification.Branch)
	message.WriteString("#")
	message.WriteString(notification.Commit)

	log.WithFields(log.Fields{
		"message": string(message.Bytes()),
	}).Debug("Prepared message")

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, message.Bytes())

	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"message": string(message.Bytes()),
		}).Fatal("Failed to sign message")
	} else {
		log.WithFields(log.Fields{
			"message": string(message.Bytes()),
			"r":       r.String(),
			"s":       s.String(),
		}).Debug("Signed message")
	}

	notification.Signature.R = r.String()
	notification.Signature.S = s.String()

	return notification
}

func (notification Notification) Verify(publicKeyPath string) bool {
	publicKeyRaw, err := ioutil.ReadFile(publicKeyPath)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"path":  publicKeyPath,
		}).Fatal("Failed to read public key")
	} else {
		log.WithFields(log.Fields{
			"path": publicKeyPath,
		}).Debug("Read public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyRaw)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to parse public key")
	} else {
		log.Debug("Parsed public key")
	}

	var r big.Int
	r.SetString(notification.Signature.R, 10)

	var s big.Int
	s.SetString(notification.Signature.S, 10)

	log.WithFields(log.Fields{
		"r": r.String(),
		"s": s.String(),
	}).Debug("Prepared message signature")

	var message bytes.Buffer
	message.WriteString(notification.Owner)
	message.WriteString("/")
	message.WriteString(notification.Repository)
	message.WriteString("@")
	message.WriteString(notification.Branch)
	message.WriteString("#")
	message.WriteString(notification.Commit)

	log.WithFields(log.Fields{
		"message": string(message.Bytes()),
	}).Debug("Prepared message")

	var authentic bool

	switch publicKey := publicKey.(type) {
	case *ecdsa.PublicKey:
		log.Debug("Public key is of type ECDSA")
		authentic = ecdsa.Verify(publicKey, message.Bytes(), &r, &s)
	default:
		log.Warn("Public key is not of type ECDSA")
		authentic = false
	}

	log.WithFields(log.Fields{
		"authentic": authentic,
	}).Debug("Verified message signature")

	return authentic
}

func Build(announcement travis.Announcement) Notification {
	log.WithFields(log.Fields{
		"repository": announcement.Payload.Repository.Name,
		"owner":      announcement.Payload.Repository.Owner,
		"commit":     announcement.Payload.Commit,
		"branch":     announcement.Payload.Branch,
	}).Debug("Building deployment notification")

	notification := Notification{}
	notification.Repository = announcement.Payload.Repository.Name
	notification.Owner = announcement.Payload.Repository.Owner
	notification.Commit = announcement.Payload.Commit
	notification.Branch = announcement.Payload.Branch

	return notification
}
