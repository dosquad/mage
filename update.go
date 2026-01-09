package mage

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper"
	"github.com/dosquad/mage/helper/build"
	"github.com/dosquad/mage/helper/must"
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/helper/web"
	"github.com/dosquad/mage/loga"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/na4ma4/go-permbits"
	"github.com/princjef/mageutil/shellcmd"
)

const (
	golangciLintConfigURL = "https://raw.githubusercontent.com/dosquad/mage/refs/heads/main/.golangci.yml"
)

// Update namespace is defined to group Update functions.
type Update mg.Namespace

// GoWorkspace create the go.work file if it is missing.
func (Update) GoWorkspace(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Update)

	goworkspaceFile := paths.MustGetGitTopLevel("go.work")

	if _, err := os.Stat(goworkspaceFile); os.IsNotExist(err) {
		return shellcmd.RunAll(
			"go work init",
			"go work use . ./magefiles",
		)
	}

	return nil
}

// GolangciLint updates the .golangci.yml from the gist.
func (Update) GolangciLint(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Update)

	golangciLintFile := paths.MustGetGitTopLevel(".golangci.yml")
	golangciLocalFile := paths.MustGetGitTopLevel(".golangci.local.yml")
	golangciRemoteFile := paths.MustGetGitTopLevel(".golangci.remote.yml")
	etag := must.Must[web.ETag](web.ETagLoadConfig())

	if !paths.FileExists(golangciLocalFile) {
		loga.PrintDebugf("Downloading file directly")
		if err := web.HTTPWriteFile(
			golangciLintConfigURL,
			golangciRemoteFile,
			nil,
			0,
		); err != nil {
			return fmt.Errorf("unable to retrieve HTTP source file: %w", err)
		}
		if paths.FileChanged(golangciLintFile, golangciRemoteFile) {
			loga.PrintFileUpdatef("Updating .golangci.yml from remote")
		}
		if err := paths.FileCopy(golangciRemoteFile, golangciLintFile, true); err != nil {
			return fmt.Errorf("unable to copy file: %w", err)
		}

		return sh.Rm(golangciRemoteFile)
	}

	loga.PrintDebugf("Downloading remote config to .golangci.remote.yml")
	if err := web.HTTPWriteFile(
		golangciLintConfigURL,
		golangciRemoteFile,
		etag.GetItem(".golangci.yml"),
		0,
	); err != nil {
		return fmt.Errorf("unable to retrieve HTTP source file: %w", err)
	}
	defer func() {
		loga.PrintDebugf("Removing remote config cache .golangci.remote.yml")
		_ = os.Remove(golangciRemoteFile)
	}()

	var yamlData []byte
	{
		var err error
		loga.PrintDebugf("Merging .golangci.remote.yml and .golangci.local.yml")
		yamlData, err = helper.MergeYaml(golangciRemoteFile, golangciLocalFile)
		if err != nil {
			return fmt.Errorf("unable to merge remote and local config: %w", err)
		}
	}

	loga.PrintDebugf("Writing merged config to .golangci.yml")
	if err := os.WriteFile(golangciLintFile, yamlData, permbits.MustString("ug=rw,o=r")); err != nil {
		return fmt.Errorf("unable to write golangci ling config: %w", err)
	}

	return nil
}

// GitIgnore updates the .gitignore from a set list.
func (Update) GitIgnore(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Update)

	gitignoreFile := paths.MustGetGitTopLevel(".gitignore")

	once := sync.Once{}

	for _, path := range []string{
		"/.makefiles",
		"/artifacts",
		"/dist",
	} {
		if !paths.FileLineExists(gitignoreFile, path) {
			once.Do(func() {
				if rn, err := paths.FileLastRune(gitignoreFile); err == nil && rn != '\n' {
					_ = paths.FileAppendLine(gitignoreFile, 0, "")
				}
			})
			loga.PrintFileUpdatef("Adding path to .gitignore: %s", path)
			if err := paths.FileAppendLine(gitignoreFile, 0, path); err != nil {
				return err
			}
		}
	}

	return nil
}

// DockerIgnore writes the .dockerignore file if it does not exist.
func (Update) DockerIgnore(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Update)

	var dcfg *build.DockerConfig
	{
		var err error
		dcfg, err = build.DockerLoadConfig()
		must.PanicIfError(err, "unable to load docker config")
	}

	dockerignoreFile := paths.MustGetGitTopLevel(".dockerignore")

	once := sync.Once{}

	for _, line := range dcfg.Ignore {
		if !paths.FileLineExists(dockerignoreFile, line) {
			once.Do(func() {
				if rn, err := paths.FileLastRune(dockerignoreFile); err == nil && rn != '\n' {
					_ = paths.FileAppendLine(dockerignoreFile, 0, "")
				}
			})
			loga.PrintFileUpdatef("Adding path to .dockerignore: %s", line)
			if err := paths.FileAppendLine(dockerignoreFile, 0, line); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateE(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Update)

	mg.CtxDeps(ctx, Update.GolangciLint)
	mg.CtxDeps(ctx, Update.GitIgnore)
	if paths.FileExists("Dockerfile") {
		mg.CtxDeps(ctx, Update.DockerIgnore)
	}

	return nil
}
