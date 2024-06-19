package mage

import (
	"context"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
)

// Test runs all tests depending on presence of files (go.mod=golang:test, etc).
func Test(ctx context.Context) error {
	if helper.FileExistsInPath("*.proto", helper.MustGetWD()) {
		mg.SerialCtxDeps(ctx, Protobuf.installProtoc)
		mg.SerialCtxDeps(ctx, Protobuf.installProtocGenGo)
		mg.SerialCtxDeps(ctx, Protobuf.installProtocGenGoGRPC)
		mg.SerialCtxDeps(ctx, Protobuf.GenGo)
		mg.SerialCtxDeps(ctx, Protobuf.GenGoGRPC)
	}

	if helper.FileExists(
		helper.MustGetWD(".golangci.yml"),
		helper.MustGetWD(".golangci.yaml"),
		helper.MustGetWD(".golangci.toml"),
		helper.MustGetWD(".golangci.json"),
	) {
		mg.SerialCtxDeps(ctx, Golang.Lint)
	}

	if helper.FileExists(
		helper.MustGetWD("go.mod"),
	) {
		mg.SerialCtxDeps(ctx, Golang.Test)
	}

	if helper.FileExists(
		helper.MustGetWD(".goreleaser.yml"),
	) {
		mg.SerialCtxDeps(ctx, Goreleaser.Healthcheck)
		mg.SerialCtxDeps(ctx, Goreleaser.Lint)
	}

	if helper.FileExists(
		helper.MustGetWD("Dockerfile"),
	) {
		mg.SerialCtxDeps(ctx, Docker.Test)
	}

	return nil
}
