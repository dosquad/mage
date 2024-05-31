package mage

import (
	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/sh"
	"go.uber.org/multierr"
)

// Clean remove generated files.
func Clean() error {
	if err := helper.BinGolangCILint().Ensure(); err != nil {
		return err
	}

	if err := helper.BinGolangCILint().Command("cache clean").Run(); err != nil {
		return err
	}

	rmFunc := func(path string) error {
		helper.PrintInfo("Removing path: %s", path)
		return sh.Rm(path)
	}

	return multierr.Combine(
		rmFunc("artifacts"),
		rmFunc(".makefiles"),
	)
}

// CleanLight avoids removing `artifacts/data` but flushes golangci-lint cache.
func CleanLight() error {
	if err := helper.BinGolangCILint().Ensure(); err != nil {
		return err
	}

	if err := helper.BinGolangCILint().Command("cache clean").Run(); err != nil {
		return err
	}

	rmFunc := func(path string) error {
		helper.PrintInfo("Removing path: %s", path)
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
