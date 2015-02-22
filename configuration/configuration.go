package configuration

import (
	"log"

	"github.com/FogCreek/mini"
)

type Configuration struct {
	Host   HostConfiguration
	Web    WebConfiguration
	Travis TravisConfiguration
}

type HostConfiguration struct {
	TTL int64
}

type WebConfiguration struct {
	Address string
	Status  bool
}

type TravisConfiguration struct {
	Authenticate bool
	Token        string
}

var (
	Conf Configuration
)

func Process() {
	config, err := mini.LoadConfiguration("iago.ini")

	if err != nil {
		log.Fatal(err)
	}

	Conf.Host.TTL = config.IntegerFromSection("Host", "TTL", 30)

	Conf.Web.Address = config.StringFromSection("Web", "Address", "127.0.01:8000")
	Conf.Web.Status = config.BooleanFromSection("Web", "Status", true)

	Conf.Travis.Authenticate = config.BooleanFromSection("Travis", "Authenticate", false)
	Conf.Travis.Token = config.StringFromSection("Travis", "Token", "")
}
