package projects

import (
	"bytes"
)

type Project struct {
	Name       string
	Owner      string
	Repository string
	Version    string
	Identifier string
	Domain     string
	Subdomain  string
	Type       string
}

func (p Project) Path() string {
	var buffer bytes.Buffer

	if p.Type == "static" {
		buffer.WriteString("/srv/")
		buffer.WriteString(p.Domain)
		buffer.WriteString("/")
		buffer.WriteString(p.Subdomain)
		buffer.WriteString("/")
	}

	return string(buffer.Bytes())
}

func (p Project) ArchivePath() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.Path())
	buffer.WriteString(".milou/")
	buffer.WriteString(p.Version)
	buffer.WriteString(".zip")

	return string(buffer.Bytes())
}

func (p Project) ArchiveLocation() string {
	var buffer bytes.Buffer

	buffer.WriteString("https://s3.amazonaws.com/")
	buffer.WriteString(p.Identifier)
	buffer.WriteString("/")
	buffer.WriteString(p.Version)
	buffer.WriteString(".zip")

	return string(buffer.Bytes())
}
