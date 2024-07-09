package cli

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	templateDir = "template"
	templateExt = ".templ"
)

var (
	ErrTemplateDir   = errors.New("could not read template directory")
	ErrNilFileReader = errors.New("file reader is nil")
)

type FileReader interface {
	ReadFile(name string) ([]byte, error)
	ReadDir(name string) ([]fs.DirEntry, error)
}

type Templater struct {
	FileReader
}

type details struct {
	Module,
	GOVersion string
}

func (t *Templater) NewProject(name, path, goversion string) error {
	if t.FileReader == nil {
		return ErrNilFileReader
	}
	templates, err := t.ReadDir(templateDir)
	if err != nil {
		return errors.Join(ErrTemplateDir, err)
	}
	projectDetails := details{
		Module:    name,
		GOVersion: goversion,
	}
	for _, template := range templates {
		name := template.Name()
		destFile := t.createDestinationFile(name, path)
		defer destFile.Close()
		content := t.readTemplate(name)
		content = t.bindDetails(content, projectDetails)
		_, _ = destFile.Write(content)
	}
	return nil
}

func (t *Templater) bindDetails(content []byte, d details) []byte {
	tmpl, _ := template.New(d.Module).Parse(string(content))
	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, d)
	return buf.Bytes()
}

func (t *Templater) readTemplate(name string) []byte {
	fileName := filepath.Join(templateDir, name)
	content, _ := t.ReadFile(fileName)
	return content
}

func (t *Templater) createDestinationFile(name, path string) *os.File {
	fileNameWithoutTemplateExt := strings.Replace(name, templateExt, "", 1)
	path = filepath.Join(path, fileNameWithoutTemplateExt)
	destFile, _ := os.Create(path)
	return destFile
}
