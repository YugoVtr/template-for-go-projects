package cli_test

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yugovtr/template-for-go-projects/cli"
)

func TestSetup(t *testing.T) {
	type args struct {
		name, project, path, template string
		assert                        func(t *testing.T, gotErr error, testCase args)
	}
	testCases := []args{
		{
			name:     "success",
			project:  GenerateRandomProjectName(t, 8),
			path:     t.TempDir(),
			template: "../.",
			assert: func(t *testing.T, gotErr error, args args) {
				t.Helper()
				assert.NoError(t, gotErr)
				AssertProject(t, args.path, args.project)
			},
		},
		{
			name:     "path does not exist",
			project:  GenerateRandomProjectName(t, 8),
			path:     "",
			template: "../.",
			assert: func(t *testing.T, gotErr error, _ args) {
				t.Helper()
				assert.ErrorIs(t, gotErr, cli.ErrNotADirectory)
			},
		},
		{
			name:     "template does not exist",
			project:  GenerateRandomProjectName(t, 8),
			path:     t.TempDir(),
			template: "",
			assert: func(t *testing.T, gotErr error, _ args) {
				t.Helper()
				assert.ErrorIs(t, gotErr, cli.ErrNotADirectory)
			},
		},
		{
			name:     "project name is empty",
			project:  "",
			path:     t.TempDir(),
			template: "../.",
			assert: func(t *testing.T, gotErr error, _ args) {
				t.Helper()
				assert.ErrorIs(t, gotErr, cli.ErrProjectNameEmpty)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := cli.Setup(testCase.project, testCase.path, testCase.template)
			testCase.assert(t, err, testCase)
		})
	}
}

func GenerateRandomProjectName(t *testing.T, size int) string {
	t.Helper()
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		t.Fatalf("could not generate random project name: %v", err)
	}
	return hex.EncodeToString(bytes)
}

func AssertProject(t *testing.T, path, project string) {
	t.Helper()
	expectedFiles := []string{
		".gitignore",
		".golangci.yml",
		"go.mod",
		"main.go",
		"Makefile",
		"README.md",
		project + ".code-workspace",
		".vscode/extensions.json",
	}
	for _, expected := range expectedFiles {
		assert.FileExists(t, fmt.Sprintf("%s/%s", path, expected))
	}
	for ignore := range cli.Ignore {
		assert.NoFileExists(t, fmt.Sprintf("%s/%s", path, ignore))
	}

	file, err := os.Open(fmt.Sprintf("%s/go.mod", path))
	require.NoError(t, err)
	t.Cleanup(func() {
		file.Close()
	})

	scanner := bufio.NewScanner(file)
	scanner.Scan() // read the first line
	assert.Equal(t, fmt.Sprintf("module %s", project), scanner.Text())
}
