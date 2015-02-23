package projects

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"net/http"
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

	buffer.WriteString(p.TemporaryPath())
	buffer.WriteString(p.Version)
	buffer.WriteString(".zip")

	return string(buffer.Bytes())
}

func (p Project) TemporaryPath() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.Path())
	buffer.WriteString(".milou/")

	return string(buffer.Bytes())
}

func (p Project) ExtractPath() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.TemporaryPath())
	buffer.WriteString(p.Subdomain)

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

func (p Project) Extract() error {
	r, err := zip.OpenReader(p.ArchivePath())
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(p.ExtractPath(), f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p Project) Download() {
	response, err := http.Get(p.ArchiveLocation())
	defer response.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	archive, err := os.Create(p.ArchivePath())
	defer archive.close()

	_, err = io.Copy(archive, response.Body)

	if err != nil {
		log.Fatal(err)
	}
}
