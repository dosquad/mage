package mage

import (
	"os"
	"path/filepath"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
	"github.com/princjef/mageutil/shellcmd"
)

const (
	commandTemplate = `CGO_ENABLED={{.CGO}} ` +
		`GOOS={{.GoOS}} ` +
		`GOARCH={{.GoArch}} ` +
		`GOARM={{.GoArm}} ` +
		`go build ` +
		`{{if .Debug}}-tags=debug {{end}}` +
		`-buildmode=default ` +
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

// Debug create debug artifact.
func (Build) Debug() error {
	paths := helper.MustCommandPaths()

	for _, cmdPath := range paths {
		ct := helper.NewCommandTemplate(true, cmdPath)
		if err := buildArtifact(ct); err != nil {
			return err
		}
	}

	return nil
}

// Release create debug artifact.
func (Build) Release() error {
	paths := helper.MustCommandPaths()

	for _, cmdPath := range paths {
		ct := helper.NewCommandTemplate(false, cmdPath)
		if err := buildArtifact(ct); err != nil {
			return err
		}
	}

	return nil
}

// ReleaseCommand create debug artifact for the specified command.
func (Build) ReleaseCommand(cmd string) error {
	paths := helper.MustCommandPaths()

	for _, cmdPath := range paths {
		ct := helper.NewCommandTemplate(false, cmdPath)
		if err := buildArtifact(ct); err != nil {
			return err
		}
	}

	return nil
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
