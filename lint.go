package mage

import (
	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
)

// Lint namespace is defined to group Lint functions.
type Lint mg.Namespace

// // Lint run all linters.
// func Lint(ctx context.Context) {
// 	mg.CtxDeps(ctx, LintGolangci)
// }

// Golangci run golangci-lint.
func (Lint) Golangci() error {
	if err := helper.BinGolangCILint().Ensure(); err != nil {
		return err
	}

	return helper.BinGolangCILint().Command("run ./... --sort-results --max-same-issues 0 --max-issues-per-linter 0").Run()
}
