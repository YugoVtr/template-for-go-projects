package main

import (
	"flag"
	"log/slog"

	"github.com/yugovtr/template-for-go-projects/cli"
)

// go build -o bin/tools/cli ./cli/cmd/main.go.
// bin/tools/cli -projectName=torugo -projectPath=/tmp/2799675634 -templateDir=.
func main() {
	projectName := flag.String("projectName", "", "Name of the project")
	projectPath := flag.String("projectPath", "", "Path to the project")
	templateDir := flag.String("templateDir", "", "Path to the template directory")
	flag.Parse()

	if err := cli.Setup(*projectName, *projectPath, *templateDir); err != nil {
		slog.Error("setup failed", "error", err)
	}
}
