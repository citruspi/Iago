package projects

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
