package configuration

import (
	"flag"
	"log"
	"os"

	"github.com/FogCreek/mini"
)

type iagoConfiguration struct {
	Hostname string
	Protocol string
	Port     int64
	Path     string
}

type checkinConfiguration struct {
	Hostname string
	Protocol string
	Port     int64
	Path     string
	TTL      int64
}

type webConfiguration struct {
	Address string
}

var (
	Iago    iagoConfiguration
	CheckIn checkinConfiguration
	Web     webConfiguration
)

func Process() {
	path := flag.String("config", "/etc/miloud.ini", "Configuration file path")
	flag.Parse()

	config, err := mini.LoadConfiguration(*path)

	if err != nil {
		log.Fatal(err)
	}

	Iago.Hostname = config.StringFromSection("Iago", "Hostname", "localhost")
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

	CheckIn.Hostname = config.StringFromSection("CheckIn", "Hostname", "")
	CheckIn.Protocol = config.StringFromSection("CheckIn", "Protocol", "http")
	CheckIn.Path = config.StringFromSection("CheckIn", "Path", "/")
	CheckIn.Port = config.IntegerFromSection("CheckIn", "Port", 0)
	CheckIn.TTL = config.IntegerFromSection("CheckIn", "TTL", 30)

	if CheckIn.Hostname == "" {
		hostname, err := os.Hostname()

		if err != nil {
			log.Fatal(err)
		}

		CheckIn.Hostname = hostname
	}

	if CheckIn.Port == 0 {
		if CheckIn.Protocol == "http" {
			CheckIn.Port = 80
		} else if CheckIn.Protocol == "https" {
			CheckIn.Port = 443
		}
	}

	Web.Address = config.StringFromSection("Web", "Address", ":9090")
}
