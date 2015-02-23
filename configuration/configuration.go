package configuration

import (
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/FogCreek/mini"
	"github.com/citruspi/Milou/projects"
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

type notificationConfiguration struct {
	Signed    bool
	PublicKey string
}

var (
	Iago         iagoConfiguration
	CheckIn      checkinConfiguration
	Web          webConfiguration
	Notification notificationConfiguration
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

	Notification.Signed = config.BooleanFromSection("Notification", "Signed", false)
	Notification.PublicKey = config.StringFromSection("Notification", "PublicKey", "/etc/iagod/key.pub")

	projectList := config.StringsFromSection("Milou", "Projects")

	for _, projectName := range projectList {
		project := projects.Project{}

		var buffer bytes.Buffer
		buffer.WriteString("Project-")
		buffer.WriteString(projectName)

		section := string(buffer.Bytes())

		project.Name = config.StringFromSection(section, "Name", "")
		project.Owner = config.StringFromSection(section, "Owner", "")
		project.Repository = config.StringFromSection(section, "Repository", "")
		project.Version = config.StringFromSection(section, "Version", "")
		project.Identifier = config.StringFromSection(section, "Identifier", "")
		project.Domain = config.StringFromSection(section, "Domain", "")
		project.Subdomain = config.StringFromSection(section, "Subdomain", "")
		project.Type = config.StringFromSection(section, "Type", "")

		projects.List = append(projects.List, project)
	}
}
