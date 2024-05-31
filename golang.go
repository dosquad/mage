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
func (Golang) InstallGovulncheck(_ context.Context) error {
	return helper.BinGovulncheck().Ensure()
}

// Vulncheck runs govulncheck.
func (Golang) Vulncheck(ctx context.Context) error {
	mg.CtxDeps(ctx, Golang.InstallGovulncheck)

	return helper.BinGovulncheck().Command("./...").Run()
}

// Test run test suite and save coverage report.
func (Golang) Test() error {
	coverPath := helper.MustGetWD("artifacts", "coverage")

	helper.MustMakeDir(coverPath, 0)

	cmd := fmt.Sprintf(""+
		"go test "+
		"-race "+
		"-covermode=atomic "+
		"-coverprofile=\"%s/cover.out\" "+
		"\"./...\"",
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
