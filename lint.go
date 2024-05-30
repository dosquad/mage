package mage

import (
	"context"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
)

// Lint run all linters.
func Lint(ctx context.Context) {
	mg.CtxDeps(ctx, LintGolangci)
}

// LintGolangci run golangci-lint.
func LintGolangci() error {
	if err := helper.BinGolangCILint().Ensure(); err != nil {
		return err
	}

	return helper.BinGolangCILint().Command("run ./... --sort-results --max-same-issues 0 --max-issues-per-linter 0").Run()
}
