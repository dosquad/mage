package mage

import (
	"context"
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

// MustGetHomeDir

// Install namespace is defined to group Install functions.
type Install mg.Namespace

// CommandAll installs a release version of a supplied command.
func (Install) CommandAll(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Install)

	pathList := paths.MustCommandPaths()

	for _, cmdPath := range pathList {
		ct := helper.NewCommandTemplate(false, cmdPath)
		if err := buildArtifact(ctx, ct); err != nil {
			return err
		}

		installDir := envs.GetEnv("INSTALL_DIR", paths.MustGetHomeDir("go", "bin"))

		if err := shellcmd.Command(fmt.Sprintf(
			`install "%s" "%s"`, ct.OutputArtifact, installDir,
		)).Run(); err != nil {
			return err
		}
	}

	return nil
}

// Command installs a release version of a supplied command.
func (Install) Command(ctx context.Context, cmd string) error {
	dyndep.CtxDeps(ctx, dyndep.Install)

	ct := helper.NewCommandTemplate(false, "./cmd/"+cmd)
	if err := buildArtifact(ctx, ct); err != nil {
		return err
	}

	installDir := envs.GetEnv("INSTALL_DIR", paths.MustGetHomeDir("go", "bin"))

	return shellcmd.Command(fmt.Sprintf(
		`install "%s" "%s"`,
		ct.OutputArtifact,
		installDir,
	)).Run()
}

// InstallE installs with release tags and the supplied arguments the command specified by INSTALL_CMD
// or the first found command if the environment is not specified.
func InstallE(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Install)

	cmdName := envs.GetEnv("INSTALL_CMD", must.Must[string](build.FirstCommandName()))

	cmdDep := mg.F(Install.Command, cmdName)
	mg.CtxDeps(ctx, cmdDep)

	return nil
}
