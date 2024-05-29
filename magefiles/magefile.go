//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"

	//mage:import
	"github.com/dosquad/mage"
)

// Install protoc, lint, test & build debug.
func Install(ctx context.Context) {
	mg.CtxDeps(ctx, mage.Lint)
	mg.CtxDeps(ctx, mage.Test)
}

var Default = Install
