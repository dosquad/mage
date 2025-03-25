package mage

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper"
	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/helper/envs"
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/loga"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
	"github.com/princjef/mageutil/shellcmd"
)

// Golang namespace is defined to group Golang functions.
type Golang mg.Namespace

// InstallGovulncheck installs govulncheck.
func (Golang) installGovulncheck(_ context.Context) error {
	return bins.Govulncheck().Ensure()
}

// Vulncheck runs govulncheck.
func (Golang) Vulncheck(ctx context.Context) error {
	mg.CtxDeps(ctx, Golang.installGovulncheck)

	return bins.Govulncheck().Command("./...").Run()
}

// Test run test suite and save coverage report.
func (Golang) Test(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Golang)
	dyndep.CtxDeps(ctx, dyndep.Test)

	coverPath := paths.MustGetArtifactPath("coverage")

	paths.MustMakeDir(coverPath, 0)

	raceArg := ""
	if v := envs.GoEnv("CGO_ENABLED", "0"); v == "1" {
		raceArg = "-race"
	}

	cmd := fmt.Sprintf(""+
		"go test "+
		"%s "+
		"-covermode=atomic "+
		"-coverprofile=\"%s/cover.out\" "+
		"\"./...\"",
		raceArg,
		coverPath)

	if err := shellcmd.Command(cmd).Run(); err != nil {
		return err
	}

	return helper.FilterCoverageOutput(filepath.Join(coverPath, "cover.out"))
}

// InstallVGT installs vgt.
func (Golang) installVGT(_ context.Context) error {
	return bins.VGT().Ensure()
}

// VisualTest runs the test suite and then renders the result using vgt (https://github.com/roblaszczak/vgt).
func (Golang) VisualTest(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Golang)
	dyndep.CtxDeps(ctx, dyndep.Test)

	mg.CtxDeps(ctx, Golang.installVGT)

	raceArg := ""
	if v := envs.GoEnv("CGO_ENABLED", "0"); v == "1" {
		raceArg = "-race"
	}

	cmd := fmt.Sprintf(""+
		"go test -count=1 -json "+
		"%s "+
		"\"./...\"",
		raceArg)

	visualTestPath := paths.MustGetArtifactPath("tests")
	paths.MustMakeDir(
		visualTestPath,
		permbits.MustString("u=rwx,go=rx"),
	)
	vgtFileName := filepath.Join(visualTestPath, "vgt-output.json")

	var output []byte
	{
		var err error
		loga.PrintCommandAlways("`%s` writing to %s", cmd, vgtFileName)
		output, err = shellcmd.Command(cmd).Output()
		if err != nil {
			return err
		}
	}

	if err := paths.FileWrite(output, vgtFileName); err != nil {
		return err
	}

	return bins.VGT().Command(
		"-dont-pass-output -from-file " + vgtFileName,
	).Run()
}

// Generate runs go generate.
func (Golang) Generate(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Golang)

	return shellcmd.Command(`go generate ./...`).Run()
}

// Lint run golangci-lint.
func (Golang) Lint(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Golang)
	dyndep.CtxDeps(ctx, dyndep.Lint)

	if err := bins.GolangCILint().Ensure(); err != nil {
		return err
	}

	return bins.GolangCILint().Command("run ./... --max-same-issues 0 --max-issues-per-linter 0").Run()
}

// Fmt run go fmt.
func (Golang) Fmt() error {
	return shellcmd.Command(`go fmt ./...`).Run()
}

// Vet run go vet.
func (Golang) Vet() error {
	return shellcmd.Command(`go vet ./...`).Run()
}
