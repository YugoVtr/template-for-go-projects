package cli_test

import (
	"crypto/rand"
	"encoding/hex"
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

func TestNewProjectFromTemplate(t *testing.T) {
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
	}
	for _, testCase := range testCases {
		t.Run(testCase.when, func(t *testing.T) {
			_, err := cli.NewProjectFromTemplate(testCase.fileReader)
			testCase.assert(t, err, testCase)
		})
	}
}

func TestTemplater_Generate(t *testing.T) {
	test := T{t}
	templater, err := cli.NewProjectFromTemplate(FileReader(".."))
	require.NoError(t, err)

	name, path := test.GenerateRandomName(8), test.TempDir()
	templater.Generate(name, path, "1.22")
	assert.Nil(t, err)
	AssertProject(test.T, name, path)
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
