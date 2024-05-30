package mage

import (
	"context"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
)

func InstallGovulncheck(_ context.Context) error {
	return helper.BinGovulncheck().Ensure()
}

func GolangVulncheck(ctx context.Context) error {
	mg.CtxDeps(ctx, InstallGovulncheck)

	return helper.BinGovulncheck().Command("./...").Run()
}
