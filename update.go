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

// Update any updates should be added here.
func Update(ctx context.Context) {
	mg.CtxDeps(ctx, UpdateGoWorkspace)
	mg.CtxDeps(ctx, UpdateGolangciLint)
	mg.CtxDeps(ctx, UpdateGitIgnore)
}

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
