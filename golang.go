package mage

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/shellcmd"
)

// Golang namespace is defined to group Golang functions.
type Golang mg.Namespace

// InstallGovulncheck installs govulncheck.
func (Golang) installGovulncheck(_ context.Context) error {
	return helper.BinGovulncheck().Ensure()
}

// Vulncheck runs govulncheck.
func (Golang) Vulncheck(ctx context.Context) error {
	mg.CtxDeps(ctx, Golang.installGovulncheck)

	return helper.BinGovulncheck().Command("./...").Run()
}

// Test run test suite and save coverage report.
func (Golang) Test() error {
	coverPath := helper.MustGetArtifactPath("coverage")

	helper.MustMakeDir(coverPath, 0)

	raceArg := ""
	if v := helper.GoEnv("CGO_ENABLED", "0"); v == "1" {
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

// Lint run golangci-lint.
func (Golang) Lint() error {
	if err := helper.BinGolangCILint().Ensure(); err != nil {
		return err
	}

	return helper.BinGolangCILint().Command("run ./... --sort-results --max-same-issues 0 --max-issues-per-linter 0").Run()
}

// Fmt run go fmt.
func (Golang) Fmt() error {
	return shellcmd.Command(`go fmt ./...`).Run()
}

// Vet run go vet.
func (Golang) Vet() error {
	return shellcmd.Command(`go vet ./...`).Run()
}
