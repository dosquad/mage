package mage

import (
	"fmt"
	"path/filepath"

	"github.com/dosquad/mage/helper"
	"github.com/princjef/mageutil/shellcmd"
)

// Test run test suite and save coverage report.
func Test() error {
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
