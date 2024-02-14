package cli

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yugovtr/template-for-go-projects/cli/errwrap"
)

var Ignore = map[string]struct{}{
	"cli":     {},
	"bin":     {},
	".git":    {},
	"go.sum":  {},
	"LICENSE": {},
}

var (
	ErrNotADirectory    = errors.New("not a directory")
	ErrProjectNameEmpty = errors.New("project name is empty")
	ErrNoGoMod          = errors.New("go.mod does not exist in the template directory")
)

func Setup(projectName, projectPath, templateDir string) error {
	logger := slog.New(slog.Default().Handler())
	logger.Info("setup", "name", projectName, "path", projectPath, "template", templateDir)
	if !DirExists(projectPath) || !DirExists(templateDir) {
		return ErrNotADirectory
	}
	if len(projectName) == 0 {
		return ErrProjectNameEmpty
	}
	if _, err := os.Stat(filepath.Join(templateDir, "go.mod")); os.IsNotExist(err) {
		return errwrap.WithMessage(templateDir, ErrNoGoMod)
	}
	if err := CopyTemplateToProject(templateDir, projectPath); err != nil {
		return err
	}
	return AddNameToProject(projectName, projectPath)
}

func DirExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CopyTemplateToProject(templateDir, projectPath string) error {
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error { //nolint:wrapcheck
		if err != nil {
			return errwrap.WithMessage("walk failed", err)
		}

		if _, ok := Ignore[info.Name()]; ok {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		dest := filepath.Join(projectPath, path[len(templateDir)-1:])
		if info.IsDir() {
			if err := os.MkdirAll(dest, os.ModePerm); err != nil {
				return errwrap.WithMessage("create dir "+dest, err)
			}
			return nil
		}
		return CopyFile(path, dest)
	})
}

func CopyFile(path, dest string) error {
	sourceFile, err := os.Open(path)
	if err != nil {
		return errwrap.WithMessage("open file "+path, err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return errwrap.WithMessage("create file "+dest, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return errwrap.WithMessage(fmt.Sprintf("copy file %s to %s", path, dest), err)
	}
	return nil
}

func AddNameToProject(projectName, projectPath string) error {
	const (
		ext                   = ".code-workspace"
		templateWorkspaceName = "template-for-go-projects" + ext
	)
	source := filepath.Join(projectPath, templateWorkspaceName)
	dest := filepath.Join(projectPath, projectName+ext)
	if err := os.Rename(source, dest); err != nil {
		return errwrap.WithMessage("rename workspace", err)
	}

	modLine := fmt.Sprintf("module %s", projectName)
	modFile := filepath.Join(projectPath, "go.mod")
	command := fmt.Sprintf("sed -i '1s/.*/%s/' %s", modLine, modFile)
	if err := exec.Command("bash", "-c", command).Run(); err != nil {
		return errwrap.WithMessage("update module name", err)
	}

	return nil
}
