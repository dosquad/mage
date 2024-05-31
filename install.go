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
func (Install) Command(_ context.Context, cmd string) error {
	ct := helper.NewCommandTemplate(false, fmt.Sprintf("./cmd/%s", cmd))
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
