package mage

import (
	"context"
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

// Build builds a release binary.
func Build(ctx context.Context) {
	mg.CtxDeps(ctx, BuildDebug)
	mg.CtxDeps(ctx, BuildRelease)
}

// BuildDebug create debug artifact.
func BuildDebug() error {
	paths := helper.MustCommandPaths()

	for _, cmdPath := range paths {
		ct := helper.NewCommandTemplate(true, cmdPath)
		var out string
		{
			var err error
			out, err = ct.Render(commandTemplate)
			if err != nil {
				return err
			}
		}

		if err := os.MkdirAll(filepath.Base(ct.OutputArtifact), permbits.MustString("a=rx,u+w")); err != nil {
			return err
		}

		if err := shellcmd.Command(out).Run(); err != nil {
			return err
		}
	}

	return nil
}

// BuildRelease create debug artifact.
func BuildRelease() error {
	paths := helper.MustCommandPaths()

	for _, cmdPath := range paths {
		ct := helper.NewCommandTemplate(false, cmdPath)
		var out string
		{
			var err error
			out, err = ct.Render(commandTemplate)
			if err != nil {
				return err
			}
		}

		if err := os.MkdirAll(filepath.Base(ct.OutputArtifact), permbits.MustString("a=rx,u+w")); err != nil {
			return err
		}

		if err := shellcmd.Command(out).Run(); err != nil {
			return err
		}
	}

	return nil
}
