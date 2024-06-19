package mage

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dosquad/mage/helper"
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
func (Build) Debug() error {
	return buildCommand(true, "")
}

// Release creates the release artifacts.
func (Build) Release() error {
	return buildCommand(false, "")
}

// DebugCommand creates a debug artifact for the specified command.
func (Build) DebugCommand(cmd string) error {
	return buildCommand(true, cmd)
}

// ReleaseCommand creates a release artifact for the specified command.
func (Build) ReleaseCommand(cmd string) error {
	return buildCommand(false, cmd)
}

func buildCommand(debug bool, cmd string) error {
	paths := helper.MustCommandPaths()

	for _, cmdPath := range paths {
		if cmd != "" && filepath.Base(cmdPath) != cmd {
			continue
		}

		ct := helper.NewCommandTemplate(debug, cmdPath)
		if err := buildPlatformIterator(ct, buildArtifact); err != nil {
			return err
		}
	}

	return nil
}

func buildPlatformIterator(ct *helper.CommandTemplate, f func(*helper.CommandTemplate) error) error {
	var err error
	platforms := strings.Split(helper.GetEnv("PLATFORMS", runtime.GOOS+"/"+runtime.GOARCH), ",")

	for _, platform := range platforms {
		ctp := helper.NewCommandTemplate(ct.Debug, ct.CommandDir)
		sp := strings.Split(platform, "/")
		if len(sp) != 2 { //nolint:mnd // "os/arch"
			continue
		}
		ctp.GoOS = sp[0]
		ctp.GoArch = sp[1]
		err = multierr.Append(err, f(ctp))
	}

	return err
}

func buildArtifact(ct *helper.CommandTemplate) error {
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
