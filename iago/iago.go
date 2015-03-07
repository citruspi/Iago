package iago

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	conf "github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/projects"
)

type CheckInPayload struct {
	Hostname     string   `json:"hostname"`
	Port         int64    `json:"port"`
	Path         string   `json:"path"`
	Protocol     string   `json:"protocol"`
	Repositories []string `json:"repositories"`
}

var (
	Endpoint string
)

func buildEndpoint() {
	var buffer bytes.Buffer

	buffer.WriteString(conf.Iago.Protocol)
	buffer.WriteString("://")
	buffer.WriteString(conf.Iago.Hostname)
	buffer.WriteString(":")
	buffer.WriteString(strconv.FormatInt(conf.Iago.Port, 10))
	buffer.WriteString(conf.Iago.Path)
	buffer.WriteString("checkin/")

	Endpoint = string(buffer.Bytes())
}

func CheckIn() {
	if Endpoint == "" {
		buildEndpoint()
	}

	var repositories []string

	for _, repository := range projects.List {
		var buffer bytes.Buffer

		buffer.WriteString(repository.Owner)
		buffer.WriteString("/")
		buffer.WriteString(repository.Repository)

		repositories = append(repositories, string(buffer.Bytes()))
	}

	payload := CheckInPayload{
		Protocol:     conf.CheckIn.Protocol,
		Hostname:     conf.CheckIn.Hostname,
		Port:         conf.CheckIn.Port,
		Path:         conf.CheckIn.Path,
		Repositories: repositories,
	}

	for {
		content, _ := json.Marshal(payload)
		body := bytes.NewBuffer(content)

		req, err := http.NewRequest("POST", Endpoint, body)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		time.Sleep(time.Duration(conf.CheckIn.TTL) * time.Second)
	}
}
