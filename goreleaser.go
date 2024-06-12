package mage

import (
	"context"
	"fmt"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
)

// Goreleaser namespace is defined to group Goreleaser functions.
type Goreleaser mg.Namespace

// installGoreleaser binary
func (Goreleaser) installGoreleaser(_ context.Context) error {
	return helper.BinGoreleaser().Ensure()
}

// Build Goreleaser config
func (Goreleaser) Build(ctx context.Context) error {
	mg.CtxDeps(ctx, Goreleaser.installGoreleaser)

	args := helper.GetEnv("GORELEASER_ARGS", "")

	return helper.BinGoreleaser().Command(
		fmt.Sprintf("build %s", args),
	).Run()
}

// Lint Goreleaser config
func (Goreleaser) Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, Goreleaser.installGoreleaser)

	args := helper.GetEnv("GORELEASER_ARGS", "")

	return helper.BinGoreleaser().Command(
		fmt.Sprintf("check %s", args),
	).Run()
}

// Healthcheck Goreleaser config
func (Goreleaser) Healthcheck(ctx context.Context) error {
	mg.CtxDeps(ctx, Goreleaser.installGoreleaser)

	args := helper.GetEnv("GORELEASER_ARGS", "")

	return helper.BinGoreleaser().Command(
		fmt.Sprintf("healthcheck %s", args),
	).Run()
}
