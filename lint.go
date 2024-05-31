package mage

import (
	"context"

	"github.com/magefile/mage/mg"
)

// Lint namespace is defined to group Lint functions.
type Lint mg.Namespace

// // Lint run all linters.
// func Lint(ctx context.Context) {
// 	mg.CtxDeps(ctx, LintGolangci)
// }

// Golangci Golang linters.
func (Lint) Golang(ctx context.Context) {
	mg.CtxDeps(ctx, Golang.Lint)
}
