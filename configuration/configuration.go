package configuration

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/citruspi/milou/projects"
)

type Configuration struct {
	Mode string `json:"mode"`
	Web  struct {
		Address string `json:"address"`
	} `json:"web"`
	TravisCI struct {
		Authenticate bool   `json:"authenticate"`
		Token        string `json:"token"`
	} `json:"travis-ci"`
	Projects []struct {
		Name       string `json:"name"`
		Owner      string `json:"owner"`
		Repository string `json:"repository"`
		Version    string `json:"version"`
		Identifier string `json:"identifier"`
		Domain     string `json:"domain"`
		Subdomain  string `json:"subdomain"`
	} `json:"projects"`
}

var (
	path *string
)

func init() {
	path = flag.String("config", "/etc/milou.conf", "Configuration file path")
	flag.Parse()

	conf := Load()

	for _, project := range conf.Projects {
		projects.List = append(projects.List, projects.Project{
			Name:       project.Name,
			Owner:      project.Owner,
			Repository: project.Repository,
			Version:    project.Version,
			Identifier: project.Identifier,
			Domain:     project.Domain,
			Subdomain:  project.Subdomain,
		})
	}
}

func Load() Configuration {
	var conf Configuration

	source, err := ioutil.ReadFile(*path)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(source, &conf)

	if err != nil {
		log.Fatal(err)
	}

	return conf
}
