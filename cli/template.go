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

type Templater struct {
	fileReader       FileReader
	templatesEntries []fs.DirEntry
}

func NewProjectFromTemplate(fileReader FileReader) (t Templater, err error) {
	if fileReader == nil {
		return t, ErrNilFileReader
	}
	templates, err := fileReader.ReadDir(templateDir)
	if err != nil {
		return t, errors.Join(ErrTemplateDir, err)
	}
	return Templater{
		fileReader:       fileReader,
		templatesEntries: templates,
	}, nil
}

func (t *Templater) Generate(name, path, goversion string) error {
	for _, template := range t.templatesEntries {
		name := template.Name()
		destFile := t.createDestinationFile(name, path)
		defer destFile.Close()
		content := t.readTemplate(name)
		content = t.bindDetails(content, map[string]any{
			"Module":    name,
			"GOVersion": goversion,
		})
		_, _ = destFile.Write(content)
	}
	return nil
}

func (t *Templater) bindDetails(content []byte, details map[string]any) []byte {
	tmpl, _ := template.New("details").Parse(string(content))
	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, details)
	return buf.Bytes()
}

func (t *Templater) readTemplate(name string) []byte {
	fileName := filepath.Join(templateDir, name)
	content, _ := t.fileReader.ReadFile(fileName)
	return content
}

func (t *Templater) createDestinationFile(name, path string) *os.File {
	fileNameWithoutTemplateExt := strings.Replace(name, templateExt, "", 1)
	path = filepath.Join(path, fileNameWithoutTemplateExt)
	destFile, _ := os.Create(path)
	return destFile
}
