package mage

import (
	"context"
)

//nolint:lll // long URL
const (
	golangciLintConfigURL = "https://gist.githubusercontent.com/na4ma4/f165f6c9af35cda6b330efdcc07a9e26/raw/7a8433c1e515bd82d1865ed9070b9caff9995703/.golangci.yml"
)

// Config checks all configuration.
func Config(_ context.Context) error {
	return nil
}
