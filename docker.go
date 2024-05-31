package mage

import (
	"context"
	"fmt"
	"strings"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/shellcmd"
	"go.uber.org/multierr"
)

func dockerCommand(
	action string,
	args ...any,
) error {
	return shellcmd.Command(fmt.Sprintf(
		"docker %s %s",
		action,
		strings.Join(helper.ArgsFromAny(args), " "),
	)).Run()
}

// Docker namespace is defined to group Docker functions.
type Docker mg.Namespace

func dockerBuildArtifacts(cfg *helper.DockerConfig) error {
	paths := helper.MustCommandPaths()

	for _, cmdPath := range paths {
		ct := helper.NewCommandTemplate(false, cmdPath)
		for _, platform := range cfg.OSArch() {
			ct.GoOS = platform.OS
			ct.GoArch = platform.Arch
		}

		if err := buildArtifact(ct); err != nil {
			return err
		}
	}

	return nil
}

func dockerBuildCommand(ctx context.Context, args []string) error {
	mg.CtxDeps(ctx, Update.DockerIgnoreFile)

	cfg := helper.MustDockerLoadConfig()
	if err := dockerBuildArtifacts(cfg); err != nil {
		return err
	}

	tags := cfg.GetTags()

	if len(tags) == 0 {
		cfg.Tag = helper.SemverBumpPatch(helper.GitSemver()) + "-" + helper.GitHash()
		tags = cfg.GetTags()
	}

	var dockerErr error
	for _, tag := range tags {
		tagArg, err := cfg.ArgsTag(tag)
		if err != nil {
			helper.PrintWarning("Unable to build Docker Image: %s", err)
			continue
		}

		helper.PrintInfo("Building Docker Image[%s]: %s:%s", strings.Join(cfg.Platforms, ","), cfg.Image, tag)

		// return nil
		dockerErr = multierr.Append(dockerErr, dockerCommand(
			"buildx build",
			cfg.Args(),
			tagArg,
			args,
		))
	}

	return dockerErr
}

// Build Builds a docker image for the current platform and "loads" the
// image into the local Docker server.
func (Docker) Build(ctx context.Context) error {
	return dockerBuildCommand(ctx,
		[]string{
			"--pull",
			"--load",
			".",
		},
	)
}

// Test Builds docker images for each target platform then
// discards the result.
func (Docker) Test(ctx context.Context) error {
	return dockerBuildCommand(ctx,
		[]string{
			"--pull",
			".",
		},
	)
}

// Push Builds docker images for each target platform and pushes
// those images, and a manifest list to the registry.
func (Docker) Push(ctx context.Context) error {
	return dockerBuildCommand(ctx,
		[]string{
			"--pull",
			"--push",
			".",
		},
	)
}
