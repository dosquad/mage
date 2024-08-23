package mage

import (
	"context"
	"os"
	"runtime"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
)

// Goreleaser namespace is defined to group Goreleaser functions.
type Goreleaser mg.Namespace

// installGoreleaser binary.
func (Goreleaser) installGoreleaser(_ context.Context) error {
	return helper.BinGoreleaser().Ensure()
}

// Build Goreleaser config.
func (Goreleaser) Build(ctx context.Context) error {
	mg.CtxDeps(ctx, Goreleaser.installGoreleaser)

	os.Setenv("GOVERSION_NR", runtime.Version())

	args := helper.GetEnv("GORELEASER_ARGS", "")

	return helper.BinGoreleaser().Command("build " + args).Run()
}

// Lint Goreleaser config.
func (Goreleaser) Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, Goreleaser.installGoreleaser)

	args := helper.GetEnv("GORELEASER_ARGS", "")

	return helper.BinGoreleaser().Command("check " + args).Run()
}

// Healthcheck Goreleaser config.
func (Goreleaser) Healthcheck(ctx context.Context) error {
	mg.CtxDeps(ctx, Goreleaser.installGoreleaser)

	args := helper.GetEnv("GORELEASER_ARGS", "")

	return helper.BinGoreleaser().Command("healthcheck " + args).Run()
}
