package mage

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/shellcmd"
)

func RunDebug(ctx context.Context, cmd string, args string) error {
	mg.CtxDeps(ctx, BuildDebug)
	ct := helper.NewCommandTemplate(true, fmt.Sprintf("./cmd/%s", cmd))

	return shellcmd.Command(fmt.Sprintf("%s %s", ct.OutputArtifact, args)).Run()
}

func RunRelease(ctx context.Context, cmd string, args string) error {
	mg.CtxDeps(ctx, BuildRelease)
	ct := helper.NewCommandTemplate(false, fmt.Sprintf("./cmd/%s", cmd))

	return shellcmd.Command(fmt.Sprintf("%s %s", ct.OutputArtifact, args)).Run()
}

func Run(ctx context.Context, args string) error {
	mg.CtxDeps(ctx, BuildDebug)

	paths := helper.MustCommandPaths()
	if len(paths) < 1 {
		return errors.New("command not found")
	}

	return RunDebug(ctx, filepath.Base(paths[0]), args)
}
