package mage

import (
	"context"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
)

// Test runs all tests depending on presence of files (go.mod=golang:test, etc).
func Test(ctx context.Context) error {
	if helper.FileExistsInPath("*.proto", helper.MustGetGitTopLevel()) {
		mg.SerialCtxDeps(ctx, Protobuf.installProtoc)
		mg.SerialCtxDeps(ctx, Protobuf.installProtocGenGo)
		mg.SerialCtxDeps(ctx, Protobuf.installProtocGenGoGRPC)
		mg.SerialCtxDeps(ctx, Protobuf.GenGo)
		mg.SerialCtxDeps(ctx, Protobuf.GenGoGRPC)
	}

	if helper.FileExists(
		helper.MustGetWD("testdata", "ca-config.json"),
	) {
		mg.SerialCtxDeps(ctx, CFSSL.Generate)
	}

	if helper.FileExists(
		helper.MustGetGitTopLevel(".golangci.yml"),
		helper.MustGetGitTopLevel(".golangci.yaml"),
		helper.MustGetGitTopLevel(".golangci.toml"),
		helper.MustGetGitTopLevel(".golangci.json"),
	) {
		mg.SerialCtxDeps(ctx, Golang.Lint)
	}

	if helper.FileExists(
		helper.MustGetGitTopLevel("go.mod"),
	) {
		mg.SerialCtxDeps(ctx, Golang.Test)
	}

	if helper.FileExists(
		helper.MustGetGitTopLevel(".goreleaser.yml"),
	) {
		mg.SerialCtxDeps(ctx, Goreleaser.Healthcheck)
		mg.SerialCtxDeps(ctx, Goreleaser.Lint)
	}

	if helper.FileExists(
		helper.MustGetGitTopLevel("Dockerfile"),
	) {
		mg.SerialCtxDeps(ctx, Docker.Test)
	}

	return nil
}
