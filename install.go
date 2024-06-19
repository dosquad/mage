package mage

import (
	"context"
	"fmt"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/shellcmd"
)

// MustGetHomeDir

// Install namespace is defined to group Install functions.
type Install mg.Namespace

// Command installs a release version of a supplied command.
func (Install) CommandAll(_ context.Context) error {
	paths := helper.MustCommandPaths()

	for _, cmdPath := range paths {
		ct := helper.NewCommandTemplate(false, cmdPath)
		if err := buildArtifact(ct); err != nil {
			return err
		}

		installDir := helper.GetEnv("INSTALL_DIR", helper.MustGetHomeDir("go", "bin"))

		if err := shellcmd.Command(fmt.Sprintf(
			`install "%s" "%s"`, ct.OutputArtifact, installDir,
		)).Run(); err != nil {
			return err
		}
	}

	return nil
}

// Command installs a release version of a supplied command.
func (Install) Command(_ context.Context, cmd string) error {
	ct := helper.NewCommandTemplate(false, "./cmd/"+cmd)
	if err := buildArtifact(ct); err != nil {
		return err
	}

	installDir := helper.GetEnv("INSTALL_DIR", helper.MustGetHomeDir("go", "bin"))

	return shellcmd.Command(fmt.Sprintf(
		`install "%s" "%s"`,
		ct.OutputArtifact,
		installDir,
	)).Run()
}

// InstallE installs with release tags and the supplied arguments the command specified by INSTALL_CMD
// or the first found command if the environment is not specified.
func InstallE(ctx context.Context) error {
	cmdName := helper.GetEnv("INSTALL_CMD", helper.Must[string](helper.FirstCommandName()))

	cmdDep := mg.F(Install.Command, cmdName)
	mg.CtxDeps(ctx, cmdDep)

	return nil
}
