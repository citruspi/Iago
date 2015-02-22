package configuration

import (
	"log"

	"github.com/FogCreek/mini"
)

type Configuration struct {
	Host HostConfiguration
	Web  WebConfiguration
}

type HostConfiguration struct {
	TTL int64
}

type WebConfiguration struct {
	Address string
	Status  bool
}

func Process() Configuration {
	config, err := mini.LoadConfiguration("iago.ini")

	if err != nil {
		log.Fatal(err)
	}

	conf := Configuration{}

	conf.Host.TTL = config.IntegerFromSection("Host", "TTL", 30)

	conf.Web.Address = config.StringFromSection("Web", "Address", "127.0.01:8000")
	conf.Web.Status = config.BooleanFromSection("Web", "Status", true)

	return conf
}
