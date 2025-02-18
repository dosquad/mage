package mage

import (
	"context"
	"os"
	"runtime"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/helper/envs"
	"github.com/magefile/mage/mg"
)

// Goreleaser namespace is defined to group Goreleaser functions.
type Goreleaser mg.Namespace

// installGoreleaser binary.
func (Goreleaser) installGoreleaser(_ context.Context) error {
	return bins.Goreleaser().Ensure()
}

// Build Goreleaser config.
func (Goreleaser) Build(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Build)
	dyndep.CtxDeps(ctx, dyndep.Goreleaser)

	mg.CtxDeps(ctx, Goreleaser.installGoreleaser)

	os.Setenv("GOVERSION_NR", runtime.Version())

	args := envs.GetEnv("GORELEASER_ARGS", "")

	return bins.Goreleaser().Command("build " + args).Run()
}

// Lint Goreleaser config.
func (Goreleaser) Lint(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Lint)
	dyndep.CtxDeps(ctx, dyndep.Goreleaser)

	mg.CtxDeps(ctx, Goreleaser.installGoreleaser)

	args := envs.GetEnv("GORELEASER_ARGS", "")

	return bins.Goreleaser().Command("check " + args).Run()
}

// Healthcheck Goreleaser config.
func (Goreleaser) Healthcheck(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Goreleaser)

	mg.CtxDeps(ctx, Goreleaser.installGoreleaser)

	args := envs.GetEnv("GORELEASER_ARGS", "")

	return bins.Goreleaser().Command("healthcheck " + args).Run()
}
