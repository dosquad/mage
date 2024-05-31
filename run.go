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

// Run namespace is defined to group Run functions.
type Run mg.Namespace

// Debug builds and executes the specified command and arguments with debug build flags.
func (Run) Debug(ctx context.Context, cmd string, args string) error {
	mg.CtxDeps(ctx, Build.Debug)
	ct := helper.NewCommandTemplate(true, fmt.Sprintf("./cmd/%s", cmd))

	return shellcmd.Command(fmt.Sprintf("%s %s", ct.OutputArtifact, args)).Run()
}

// Release builds and executes the specified command and arguments with release build flags.
func (Run) Release(ctx context.Context, cmd string, args string) error {
	mg.CtxDeps(ctx, Build.Release)
	ct := helper.NewCommandTemplate(false, fmt.Sprintf("./cmd/%s", cmd))

	return shellcmd.Command(fmt.Sprintf("%s %s", ct.OutputArtifact, args)).Run()
}

// Runc builds and executes the first found command with debug tags and the supplied arguments.
func Runc(ctx context.Context, args string) error {
	mg.CtxDeps(ctx, Build.Debug)

	paths := helper.MustCommandPaths()
	if len(paths) < 1 {
		return errors.New("command not found")
	}

	ct := helper.NewCommandTemplate(true, fmt.Sprintf("./cmd/%s", filepath.Base(paths[0])))

	return shellcmd.Command(fmt.Sprintf("%s %s", ct.OutputArtifact, args)).Run()
}
