package cli_test

import (
	"crypto/rand"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yugovtr/template-for-go-projects/cli"
)

type T struct {
	*testing.T
}

func (test T) GenerateRandomName(size int) string {
	test.Helper()
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		test.Fatalf("could not generate random project name: %v", err)
	}
	return hex.EncodeToString(bytes)
}

func TestTemplater_NewProject(t *testing.T) {
	test := T{t}

	type arg struct {
		when, name, path string
		fileReader       cli.FileReader
		assert           func(*testing.T, error, arg)
	}

	testCases := []arg{
		{
			when: "file reader is nil",
			assert: func(t *testing.T, err error, _ arg) {
				assert.ErrorIs(t, err, cli.ErrNilFileReader)
			},
		},
		{
			when:       "template directory does not exist",
			fileReader: FileReader(test.TempDir()),
			assert: func(t *testing.T, err error, _ arg) {
				assert.ErrorIs(t, err, cli.ErrTemplateDir)
			},
		},
		{
			when:       "create project successfully",
			fileReader: FileReader(".."),
			name:       test.GenerateRandomName(8),
			path:       test.TempDir(),
			assert: func(t *testing.T, err error, testCase arg) {
				assert.Nil(t, err)
				AssertProject(test.T, testCase.name, testCase.path)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.when, func(t *testing.T) {
			templater := cli.Templater{FileReader: testCase.fileReader}
			err := templater.NewProject(testCase.name, testCase.path, "1.22")
			testCase.assert(t, err, testCase)
		})
	}
}

func AssertProject(t *testing.T, name, path string) {
	t.Helper()
	expectedFiles := []string{
		".gitignore",
		".golangci.yml",
		"go.mod",
		"main.go",
		"README.md",
		"project.code-workspace",
	}
	for _, expected := range expectedFiles {
		fileName := filepath.Join(path, expected)
		assert.FileExists(t, fileName)

		file, err := os.ReadFile(fileName)
		require.NoError(t, err)
		assert.NotEmpty(t, file)
	}
}

type FileReader string

// Interface guard
var _ cli.FileReader = FileReader("")

func (f FileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(string(f), name))
}

func (f FileReader) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(string(f), name))
}
