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

	return multierr.Combine(
		sh.Rm("artifacts"),
		sh.Rm(".makefiles"),
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

	return multierr.Combine(
		sh.Rm(".makefiles"),
		sh.Rm("artifacts/bin"),
		sh.Rm("artifacts/build"),
		sh.Rm("artifacts/protobuf"),
	)
}
