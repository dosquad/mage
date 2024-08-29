package mage

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper"
	"github.com/dosquad/mage/loga"
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

func dockerBuildArtifacts(ctx context.Context, cfg *helper.DockerConfig) error {
	paths := helper.MustCommandPaths()

	for _, cmdPath := range paths {
		ct := helper.NewCommandTemplate(false, cmdPath)
		for _, platform := range cfg.OSArch() {
			ct.GoOS = platform.OS
			ct.GoArch = platform.Arch
			if err := buildArtifact(ctx, ct); err != nil {
				return err
			}
		}
	}

	return nil
}

func dockerBuildCommand(ctx context.Context, args []string) error {
	cfg := helper.Must[*helper.DockerConfig](helper.DockerLoadConfig())

	var once sync.Once
	tags := cfg.GetTags()

	if len(tags) == 0 {
		cfg.Tag = helper.SemverBumpPatch(helper.GitSemver()) + "-" + helper.GitHash()
		tags = cfg.GetTags()
	}

	var dockerErr error
	for _, tag := range tags {
		tagArg, err := cfg.ArgsTag(tag)
		if err != nil {
			loga.PrintWarning("Unable to build Docker Image: %s", err)
			continue
		}

		{
			var outErr error
			once.Do(func() {
				mg.CtxDeps(ctx, Update.DockerIgnore)
				outErr = dockerBuildArtifacts(ctx, cfg)
			})
			if outErr != nil {
				return outErr
			}
		}

		loga.PrintInfo("Building Docker Image[%s]: %s:%s", strings.Join(cfg.Platforms, ","), cfg.GetImage(), tag)

		// return nil
		dockerErr = multierr.Append(dockerErr, dockerCommand(
			"buildx build",
			cfg.Args(ctx),
			tagArg,
			args,
		))
	}

	return dockerErr
}

// Build Builds a docker image for the current platform and "loads" the
// image into the local Docker server.
func (Docker) Build(ctx context.Context) error {
	ctx = context.WithValue(ctx, helper.DockerLocalPlatform, true)
	dyndep.CtxDeps(ctx, dyndep.Build)
	dyndep.CtxDeps(ctx, dyndep.Docker)
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
	dyndep.CtxDeps(ctx, dyndep.Docker)
	dyndep.CtxDeps(ctx, dyndep.Test)
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
