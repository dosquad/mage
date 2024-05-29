package mage

import (
	"github.com/magefile/mage/sh"
	"go.uber.org/multierr"
)

// Clean remove generated files.
func Clean() error {
	if err := golangciLint.Command("cache clean").Run(); err != nil {
		return err
	}

	return multierr.Combine(
		sh.Rm("artifacts"),
		sh.Rm(".makefiles"),
	)
}

// CleanLight remove generated files.
func CleanLight() error {
	if err := golangciLint.Command("cache clean").Run(); err != nil {
		return err
	}

	for _, path := range []string{
		"artifacts/bin",
		"artifacts/build",
		"artifacts/protobuf",
	} {
		if err := sh.Rm(path); err != nil {
			return err
		}
	}

	return multierr.Combine(
		sh.Rm(".makefiles"),
		sh.Rm("artifacts/bin"),
		sh.Rm("artifacts/build"),
		sh.Rm("artifacts/protobuf"),
	)
}
