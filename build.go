package mage

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper"
	"github.com/dosquad/mage/helper/envs"
	"github.com/dosquad/mage/helper/paths"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
	"github.com/princjef/mageutil/shellcmd"
	"go.uber.org/multierr"
)

const (
	commandTemplate = `CGO_ENABLED={{.CGO}} ` +
		`GOOS={{.GoOS}} ` +
		`GOARCH={{.GoArch}} ` +
		`GOARM={{.GoArm}} ` +
		`go build ` +
		`{{if .Debug}}-tags=debug {{else}}-tags=release {{end}}` +
		`-buildmode=default ` +
		`{{if not .Debug}}-trimpath {{end}}` +
		`-v ` +
		`-ldflags "{{.LDFlags}}" ` +
		`-o "{{.OutputArtifact}}" ` +
		`"{{.CommandDir}}"`
)

// Build namespace is defined to group Build functions.
type Build mg.Namespace

// // Build builds a release binary.
// func Build(ctx context.Context) {
// 	mg.CtxDeps(ctx, BuildDebug)
// 	mg.CtxDeps(ctx, BuildRelease)
// }

// Debug creates the debug artifacts.
func (Build) Debug(ctx context.Context) error {
	return buildCommand(ctx, true, "")
}

// Release creates the release artifacts.
func (Build) Release(ctx context.Context) error {
	return buildCommand(ctx, false, "")
}

// DebugCommand creates a debug artifact for the specified command.
func (Build) DebugCommand(ctx context.Context, cmd string) error {
	return buildCommand(ctx, true, cmd)
}

// ReleaseCommand creates a release artifact for the specified command.
func (Build) ReleaseCommand(ctx context.Context, cmd string) error {
	return buildCommand(ctx, false, cmd)
}

func buildCommand(ctx context.Context, debug bool, cmd string) error {
	dyndep.CtxDeps(ctx, dyndep.Build)
	pathList := paths.MustCommandPaths()

	for _, cmdPath := range pathList {
		if cmd != "" && filepath.Base(cmdPath) != cmd {
			continue
		}

		ct := helper.NewCommandTemplate(debug, cmdPath)
		if err := buildPlatformIterator(ctx, ct, buildArtifact); err != nil {
			return err
		}
	}

	return nil
}

func buildPlatformIterator(
	ctx context.Context,
	ct *helper.CommandTemplate,
	f func(context.Context, *helper.CommandTemplate) error,
) error {
	var err error
	platforms := strings.Split(envs.GetEnv("PLATFORMS", runtime.GOOS+"/"+runtime.GOARCH), ",")

	for _, platform := range platforms {
		ctp := helper.NewCommandTemplate(ct.Debug, ct.CommandDir)
		ctp.AdditionalArtifacts = ct.AdditionalArtifacts
		sp := strings.Split(platform, "/")
		if len(sp) != 2 { //nolint:mnd // "os/arch"
			continue
		}
		if after, ok := strings.CutPrefix(sp[1], "armv"); ok {
			ctp.SetPlatform(sp[0], "arm", after)
		} else {
			ctp.SetPlatform(sp[0], sp[1], "")
		}
		err = multierr.Append(err, f(ctx, ctp))
	}

	return err
}

func buildArtifact(_ context.Context, ct *helper.CommandTemplate) error {
	var out string
	{
		var err error
		out, err = ct.Render(commandTemplate)
		if err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(ct.OutputArtifact), permbits.MustString("a=rx,u+w")); err != nil {
		return err
	}

	if err := shellcmd.Command(out).Run(); err != nil {
		return err
	}

	return nil
}
