package mage

import (
	"context"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/bintool"
)

// Lint run all linters.
func Lint(ctx context.Context) {
	mg.CtxDeps(ctx, LintGolang)
}

//nolint:gochecknoglobals // ignore globals
var golangciLint = bintool.Must(bintool.New(
	"golangci-lint{{.BinExt}}",
	golangciLintVersion,
	"https://github.com/golangci/golangci-lint/releases/download/"+
		"v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
	bintool.WithFolder(helper.MustGetGoBin()),
))

// LintGolang run golangci-lint.
func LintGolang() error {
	if err := golangciLint.Ensure(); err != nil {
		return err
	}

	return golangciLint.Command("run ./... --sort-results --max-same-issues 0 --max-issues-per-linter 0").Run()
}
