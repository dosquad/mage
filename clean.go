package mage

import (
	"context"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/loga"
	"github.com/magefile/mage/sh"
	"go.uber.org/multierr"
)

// Clean remove generated files.
func Clean(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Clean)

	if bins.GolangCILint().IsInstalled() {
		if err := bins.GolangCILint().Ensure(); err != nil {
			return err
		}

		if err := bins.GolangCILint().Command("cache clean").Run(); err != nil {
			return err
		}
	}

	rmFunc := func(path string) error {
		loga.PrintInfo("Removing path: %s", path)
		return sh.Rm(path)
	}

	return multierr.Combine(
		rmFunc("artifacts"),
		rmFunc(".makefiles"),
	)
}

// CleanLight avoids removing `artifacts/data` but flushes golangci-lint cache.
func CleanLight(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Clean)

	if err := bins.GolangCILint().Ensure(); err != nil {
		return err
	}

	if err := bins.GolangCILint().Command("cache clean").Run(); err != nil {
		return err
	}

	rmFunc := func(path string) error {
		loga.PrintInfo("Removing path: %s", path)
		return sh.Rm(path)
	}

	return multierr.Combine(
		rmFunc(".makefiles"),
		rmFunc("artifacts/.versioncache.yaml"),
		rmFunc("artifacts/bin"),
		rmFunc("artifacts/build"),
		rmFunc("artifacts/config"),
		rmFunc("artifacts/lint"),
		rmFunc("artifacts/protobuf"),
	)
}
