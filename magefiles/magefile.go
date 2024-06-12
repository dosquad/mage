//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"

	//mage:import
	"github.com/dosquad/mage"
)

// TestLocal protoc, lint, test & build debug.
func TestLocal(ctx context.Context) {
	mg.CtxDeps(ctx, mage.Golang.Lint)
	mg.CtxDeps(ctx, mage.Golang.Test)
}

var Default = TestLocal
