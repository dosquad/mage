package mage

import (
	"context"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
)

// Wire namespace is defined to group Wire functions.
type Wire mg.Namespace

// installWireBinary installs govulncheck.
func (Wire) installWireBinary(_ context.Context) error {
	return helper.BinWire().Ensure()
}

// Generate install and generate golang wire dependency files.
func (Wire) Generate(ctx context.Context) error {
	mg.CtxDeps(ctx, Wire.installWireBinary)

	return helper.BinWire().Command("gen ./...").Run()
}

// Generate install and generate golang wire dependency files.
func (Wire) Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, Wire.installWireBinary)

	return helper.BinWire().Command("check ./...").Run()
}
