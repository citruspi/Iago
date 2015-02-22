package iago

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	conf "github.com/citruspi/Milou/configuration"
)

type CheckInPayload struct {
	Hostname string `json:"hostname"`
	Port     int64  `json:"port"`
	Path     string `json:"path"`
	Protocol string `json:"protocol"`
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

	payload := CheckInPayload{
		Protocol: conf.CheckIn.Protocol,
		Hostname: conf.CheckIn.Hostname,
		Port:     conf.CheckIn.Port,
		Path:     conf.CheckIn.Path,
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
