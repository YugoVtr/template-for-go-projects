package main

import (
	"embed"
	"flag"
	"log/slog"

	"github.com/yugovtr/template-for-go-projects/cli"
)

//go:embed template/*
var templates embed.FS

const defaultGoVersion = "1.22.0"

func main() {
	projectName := flag.String("name", "", "Name of the project")
	projectPath := flag.String("path", "", "Path to the project")
	goversion := flag.String("goversion", defaultGoVersion, "Go version to use in the project")
	flag.Parse()

	if *projectName == "" || *projectPath == "" {
		flag.PrintDefaults()
		return
	}

	templater, err := cli.NewProjectFromTemplate(templates)
	if err != nil {
		slog.Error("could not create project from template", "error", err)
	}
	if err := templater.Generate(*projectName, *projectPath, *goversion); err != nil {
		slog.Error("could not generate project", "error", err)
	}
}
