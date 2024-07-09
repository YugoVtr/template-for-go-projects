package main

import (
	"embed"
	"flag"
	"log/slog"

	"github.com/yugovtr/template-for-go-projects/cli"
)

//go:embed template/*
var templates embed.FS

func main() {
	projectName := flag.String("name", "", "Name of the project")
	projectPath := flag.String("path", "", "Path to the project")
	goversion := flag.String("goversion", "1.22.0", "Go version to use in the project")
	flag.Parse()

	if *projectName == "" || *projectPath == "" {
		flag.PrintDefaults()
		return
	}

	templater := cli.Templater{FileReader: templates}
	if err := templater.NewProject(*projectName, *projectPath, *goversion); err != nil {
		slog.Error("setup failed", "error", err)
	}
}
