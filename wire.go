package mage

import (
	"context"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper/bins"
	"github.com/magefile/mage/mg"
)

// Wire namespace is defined to group Wire functions.
type Wire mg.Namespace

// installWireBinary installs govulncheck.
func (Wire) installWireBinary(_ context.Context) error {
	return bins.Wire().Ensure()
}

// Generate install and generate golang wire dependency files.
func (Wire) Generate(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Wire)

	mg.CtxDeps(ctx, Wire.installWireBinary)

	return bins.Wire().Command("gen ./...").Run()
}

// Lint golang wire dependency files.
func (Wire) Lint(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Lint)
	dyndep.CtxDeps(ctx, dyndep.Wire)

	mg.CtxDeps(ctx, Wire.installWireBinary)

	return bins.Wire().Command("check ./...").Run()
}
