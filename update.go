package mage

import (
	"context"
	"os"
	"sync"

	"github.com/dosquad/mage/helper"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/shellcmd"
)

//nolint:lll // long URL
const (
	golangciLintConfigURL = "https://gist.githubusercontent.com/na4ma4/f165f6c9af35cda6b330efdcc07a9e26/raw/7a8433c1e515bd82d1865ed9070b9caff9995703/.golangci.yml"
)

// Update executes the set of updates.
func Update(ctx context.Context) {
	mg.CtxDeps(ctx, UpdateGoWorkspace)
	mg.CtxDeps(ctx, UpdateGolangciLint)
	mg.CtxDeps(ctx, UpdateGitIgnore)
}

// UpdateGoWorkspace create the go.work file if it is missing.
func UpdateGoWorkspace() error {
	goworkspaceFile := helper.MustGetWD("go.work")

	if _, err := os.Stat(goworkspaceFile); os.IsNotExist(err) {
		return shellcmd.RunAll(
			"go work init",
			"go work use . ./magefiles",
		)
	}

	return nil
}

// UpdateGolangciLint updates the .golangci.yml from the gist.
func UpdateGolangciLint() error {
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

// UpdateGitIgnore updates the .gitignore from a set list.
func UpdateGitIgnore() error {
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
