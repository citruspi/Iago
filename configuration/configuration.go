package configuration

import (
	"bytes"
	"flag"
	"log"

	"github.com/FogCreek/mini"
	"github.com/citruspi/milou/projects"
)

type milouConfiguration struct {
	Mode string
}

type webConfiguration struct {
	Address string
}

type travisConfiguration struct {
	Authenticate bool
	Token        string
}

var (
	Milou  milouConfiguration
	Web    webConfiguration
	Travis travisConfiguration
)

func Process() {
	path := flag.String("config", "/etc/milou.ini", "Configuration file path")
	flag.Parse()

	config, err := mini.LoadConfiguration(*path)

	if err != nil {
		log.Fatal(err)
	}

	Milou.Mode = config.StringFromSection("Milou", "Mode", "server")

	Web.Address = config.StringFromSection("Web", "Address", "127.0.01:8000")

	Travis.Authenticate = config.BooleanFromSection("Travis", "Authenticate", false)
	Travis.Token = config.StringFromSection("Travis", "Token", "")

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
