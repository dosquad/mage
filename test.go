package mage

import (
	"context"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper/paths"
	"github.com/magefile/mage/mg"
)

// Test runs all tests depending on presence of files (go.mod=golang:test, etc).
func Test(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Test)

	if paths.FileExistsInPath("*.proto", paths.MustGetGitTopLevel()) {
		mg.SerialCtxDeps(ctx, Protobuf.installProtoc)
		mg.SerialCtxDeps(ctx, Protobuf.installProtocGenGo)
		mg.SerialCtxDeps(ctx, Protobuf.installProtocGenGoGRPC)
		mg.SerialCtxDeps(ctx, Protobuf.GenGo)
		mg.SerialCtxDeps(ctx, Protobuf.GenGoGRPC)
	}

	if paths.FileExists(
		paths.MustGetWD("testdata", "ca-config.json"),
	) {
		mg.SerialCtxDeps(ctx, CFSSL.Generate)
	}

	if paths.FileExists(
		paths.MustGetGitTopLevel(".golangci.yml"),
		paths.MustGetGitTopLevel(".golangci.yaml"),
		paths.MustGetGitTopLevel(".golangci.toml"),
		paths.MustGetGitTopLevel(".golangci.json"),
	) {
		mg.SerialCtxDeps(ctx, Golang.Lint)
	}

	if paths.FileExists(
		paths.MustGetGitTopLevel("go.mod"),
	) {
		mg.SerialCtxDeps(ctx, Golang.Test)
	}

	if paths.FileExists(
		paths.MustGetGitTopLevel(".goreleaser.yml"),
	) {
		mg.SerialCtxDeps(ctx, Goreleaser.Healthcheck)
		mg.SerialCtxDeps(ctx, Goreleaser.Lint)
	}

	if paths.FileExists(
		paths.MustGetGitTopLevel("Dockerfile"),
	) {
		mg.SerialCtxDeps(ctx, Docker.Test)
	}

	return nil
}
