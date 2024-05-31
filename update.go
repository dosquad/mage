package mage

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/dosquad/mage/helper"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
	"github.com/princjef/mageutil/shellcmd"
)

//nolint:lll // long URL
const (
	golangciLintConfigURL = "https://gist.githubusercontent.com/na4ma4/f165f6c9af35cda6b330efdcc07a9e26/raw/7a8433c1e515bd82d1865ed9070b9caff9995703/.golangci.yml"
)

// Update namespace is defined to group Update functions.
type Update mg.Namespace

// // Update executes the set of updates.
// func Update(ctx context.Context) {
// 	mg.CtxDeps(ctx, UpdateGoWorkspace)
// 	mg.CtxDeps(ctx, UpdateGolangciLint)
// 	mg.CtxDeps(ctx, UpdateGitIgnore)
// }

// GoWorkspace create the go.work file if it is missing.
func (Update) GoWorkspace() error {
	goworkspaceFile := helper.MustGetWD("go.work")

	if _, err := os.Stat(goworkspaceFile); os.IsNotExist(err) {
		return shellcmd.RunAll(
			"go work init",
			"go work use . ./magefiles",
		)
	}

	return nil
}

// GolangciLint updates the .golangci.yml from the gist.
func (Update) GolangciLint() error {
	golangciLintFile := helper.MustGetWD(".golangci.yml")

	// if _, err := os.Stat(golangciLintFile); os.IsNotExist(err) || force {
	return helper.HTTPWriteFile(
		golangciLintConfigURL,
		golangciLintFile,
		0,
	)
	// }

	// return nil
}

// GitIgnore updates the .gitignore from a set list.
func (Update) GitIgnore() error {
	gitignoreFile := helper.MustGetWD(".gitignore")

	once := sync.Once{}

	for _, path := range []string{
		"/artifacts",
	} {
		if !helper.FileLineExists(gitignoreFile, path) {
			once.Do(func() {
				if rn, err := helper.FileLastRune(gitignoreFile); err == nil && rn != '\n' {
					_ = helper.FileAppendLine(gitignoreFile, 0, "")
				}
			})
			color.Blue("Adding path to .gitignore: %s", path)
			if err := helper.FileAppendLine(gitignoreFile, 0, path); err != nil {
				return err
			}
		}
	}

	return nil
}

// DockerIgnoreFile writes the .dockerignore file if it does not exist.
func (Update) DockerIgnoreFile() error {
	dockerignoreFile := helper.MustGetWD(".gitignore")

	buf := bytes.NewBuffer(nil)
	for _, line := range []string{
		".makefiles",
		".git",
		".github",
	} {
		if _, err := buf.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("unable to write line to buffer: %w", err)
		}
	}

	if !helper.FileExists(dockerignoreFile) {
		return os.WriteFile(
			dockerignoreFile,
			buf.Bytes(),
			permbits.MustString("a=rw"),
		)
	}

	return nil
}
