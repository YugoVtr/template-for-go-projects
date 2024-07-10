package cli_test

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/yugovtr/template-for-go-projects/cli"
)

type FileReader string

// Interface guard
var _ cli.FileReader = FileReader("")

func (f FileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(string(f), name))
}

func (f FileReader) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(string(f), name))
}
