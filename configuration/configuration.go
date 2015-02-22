package configuration

import (
	"flag"
	"log"

	"github.com/FogCreek/mini"
)

type iagoConfiguration struct {
	Hostname string
	Protocol string
	Port     string
	Path     string
}

var (
	Iago iagoConfiguration
)

func Process() {
	path := flag.String("config", "/etc/miloud.ini", "Configuration file path")
	flag.Parse()

	config, err := mini.LoadConfiguration(*path)

	if err != nil {
		log.Fatal(err)
	}

	Iago.Hostname = config.StringFromSection("Iago", "Hostname", "")
	Iago.Protocol = config.StringFromSection("Iago", "Protocol", "http")
	Iago.Path = config.StringFromSection("Iago", "Path", "/")
	Iago.Port = config.IntegerFromSection("Iago", "Port", 0)

	if Iago.Port == 0 {
		if Iago.Protocol == "http" {
			Iago.Port = 80
		} else if Iago.Protocol == "https" {
			Iago.Port = 443
		}
	}
}
