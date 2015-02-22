package iago

import (
	"bytes"
	"strconv"

	conf "github.com/citruspi/Milou/configuration"
)

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
