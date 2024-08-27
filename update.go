package mage

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/na4ma4/go-permbits"
	"github.com/princjef/mageutil/shellcmd"
)

const (
	golangciLintConfigURL = "https://gist.githubusercontent.com/na4ma4/f165f6c9af35cda6b330efdcc07a9e26/raw/.golangci.yml"
)

// Update namespace is defined to group Update functions.
type Update mg.Namespace

// GoWorkspace create the go.work file if it is missing.
func (Update) GoWorkspace() error {
	goworkspaceFile := helper.MustGetGitTopLevel("go.work")

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
	golangciLintFile := helper.MustGetGitTopLevel(".golangci.yml")
	golangciLocalFile := helper.MustGetGitTopLevel(".golangci.local.yml")
	golangciRemoteFile := helper.MustGetGitTopLevel(".golangci.remote.yml")
	etag := helper.Must[helper.ETag](helper.ETagLoadConfig())

	if !helper.FileExists(golangciLocalFile) {
		helper.PrintDebug("Downloading file directly")
		if err := helper.HTTPWriteFile(
			golangciLintConfigURL,
			golangciRemoteFile,
			nil,
			0,
		); err != nil {
			return fmt.Errorf("unable to retrieve HTTP source file: %w", err)
		}
		if helper.FileChanged(golangciLintFile, golangciRemoteFile) {
			helper.PrintFileUpdate("Updating .golangci.yml from remote")
		}
		if err := helper.FileCopy(golangciRemoteFile, golangciLintFile, true); err != nil {
			return fmt.Errorf("unable to copy file: %w", err)
		}

		return sh.Rm(golangciRemoteFile)
	}

	helper.PrintDebug("Downloading remote config to .golangci.remote.yml")
	if err := helper.HTTPWriteFile(
		golangciLintConfigURL,
		golangciRemoteFile,
		etag.GetItem(".golangci.yml"),
		0,
	); err != nil {
		return fmt.Errorf("unable to retrieve HTTP source file: %w", err)
	}
	defer func() {
		helper.PrintDebug("Removing remote config cache .golangci.remote.yml")
		_ = os.Remove(golangciRemoteFile)
	}()

	var yamlData []byte
	{
		var err error
		helper.PrintDebug("Merging .golangci.remote.yml and .golangci.local.yml")
		yamlData, err = helper.MergeYaml(golangciRemoteFile, golangciLocalFile)
		if err != nil {
			return fmt.Errorf("unable to merge remote and local config: %w", err)
		}
	}

	helper.PrintDebug("Writing merged config to .golangci.yml")
	if err := os.WriteFile(golangciLintFile, yamlData, permbits.MustString("ug=rw,o=r")); err != nil {
		return fmt.Errorf("unable to write golangci ling config: %w", err)
	}

	return nil
}

// GitIgnore updates the .gitignore from a set list.
func (Update) GitIgnore() error {
	gitignoreFile := helper.MustGetGitTopLevel(".gitignore")

	once := sync.Once{}

	for _, path := range []string{
		"/.makefiles",
		"/artifacts",
		"/dist",
	} {
		if !helper.FileLineExists(gitignoreFile, path) {
			once.Do(func() {
				if rn, err := helper.FileLastRune(gitignoreFile); err == nil && rn != '\n' {
					_ = helper.FileAppendLine(gitignoreFile, 0, "")
				}
			})
			helper.PrintFileUpdate("Adding path to .gitignore: %s", path)
			if err := helper.FileAppendLine(gitignoreFile, 0, path); err != nil {
				return err
			}
		}
	}

	return nil
}

// DockerIgnore writes the .dockerignore file if it does not exist.
func (Update) DockerIgnore() error {
	var dcfg *helper.DockerConfig
	{
		var err error
		dcfg, err = helper.DockerLoadConfig()
		helper.PanicIfError(err, "unable to load docker config")
	}

	dockerignoreFile := helper.MustGetGitTopLevel(".dockerignore")

	once := sync.Once{}

	for _, line := range dcfg.Ignore {
		if !helper.FileLineExists(dockerignoreFile, line) {
			once.Do(func() {
				if rn, err := helper.FileLastRune(dockerignoreFile); err == nil && rn != '\n' {
					_ = helper.FileAppendLine(dockerignoreFile, 0, "")
				}
			})
			helper.PrintFileUpdate("Adding path to .dockerignore: %s", line)
			if err := helper.FileAppendLine(dockerignoreFile, 0, line); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateE(ctx context.Context) error {
	mg.CtxDeps(ctx, Update.GolangciLint)
	mg.CtxDeps(ctx, Update.GitIgnore)
	if helper.FileExists("Dockerfile") {
		mg.CtxDeps(ctx, Update.DockerIgnore)
	}

	return nil
}
