package mage

import (
	"context"
	"errors"
	"fmt"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper"
	"github.com/dosquad/mage/helper/build"
	"github.com/dosquad/mage/helper/envs"
	"github.com/dosquad/mage/helper/must"
	"github.com/dosquad/mage/helper/paths"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/shellcmd"
)

// Run namespace is defined to group Run functions.
type Run mg.Namespace

// Debug builds and executes the specified command and arguments with debug build flags.
func (Run) Debug(ctx context.Context, cmd string, args string) error {
	dyndep.CtxDeps(ctx, dyndep.Run)

	mg.CtxDeps(ctx, Build.Debug)
	ct := helper.NewCommandTemplate(true, "./cmd/"+cmd)

	return shellcmd.Command(ct.OutputArtifact + " " + args).Run()
}

// Release builds and executes the specified command and arguments with release build flags.
func (Run) Release(ctx context.Context, cmd string, args string) error {
	dyndep.CtxDeps(ctx, dyndep.Run)

	mg.CtxDeps(ctx, Build.Release)
	ct := helper.NewCommandTemplate(false, "./cmd/"+cmd)

	return shellcmd.Command(fmt.Sprintf("%s %s", ct.OutputArtifact, args)).Run()
}

// Mirrord start service with Mirrord intercepts.
func (Run) Mirrord(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Run)
	dyndep.CtxDeps(ctx, dyndep.Mirrord)

	cfg := must.Must[*build.DockerConfig](build.DockerLoadConfig())
	cfgFile := paths.MustGetGitTopLevel("mirrord.yaml")

	if !paths.FileExists(cfgFile) {
		return fmt.Errorf("Mirrord configuration file (%s) missing", cfgFile)
	}

	mg.CtxDeps(ctx, Build.Debug)
	ct := helper.NewCommandTemplate(true, "./cmd/"+must.Must[string](build.FirstCommandName()))

	// targetCmd := fmt.Sprintf("artifacts/build/debug/%s/%s/%s", Cfg.OOS, Cfg.Arch, Cfg.BaseDir)

	if cfg.Mirrord.Targetless {
		return shellcmd.Command(
			fmt.Sprintf(
				"mirrord exec --config-file %s %s",
				cfgFile,
				ct.OutputArtifact,
			),
		).Run()
	}

	var targetPod string
	{
		if cfg.Kubernetes.PodSelector == "" {
			panic(errors.New("kubernetes.pod-selector in .docker.yml must not be empty, " +
				"example: 'deployment=slackrobot-router'"))
		}
		targetPod = must.Must[string](
			build.KubernetesGetPodWithSelector(cfg.Kubernetes.PodSelector),
		)
	}
	return shellcmd.Command(
		fmt.Sprintf(
			"mirrord exec --config-file %s -t %s -n %s %s",
			cfgFile,
			targetPod,
			must.Must[string](build.KubernetesGetCurrentContext()),
			ct.OutputArtifact,
		),
	).Run()
}

// RunE builds with debug tags and the supplied arguments the command specified by RUN_CMD
// or the first found command if the environment is not specified.
func RunE(ctx context.Context, args string) error {
	dyndep.CtxDeps(ctx, dyndep.Run)

	cmdName := envs.GetEnv("RUN_CMD", must.Must[string](build.FirstCommandName()))
	ct := helper.NewCommandTemplate(true, "./cmd/"+cmdName)

	if err := buildArtifact(ctx, ct); err != nil {
		return err
	}

	return shellcmd.Command(fmt.Sprintf("%s %s", ct.OutputArtifact, args)).Run()
}
