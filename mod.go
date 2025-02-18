package mage

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/helper/build"
	"github.com/dosquad/mage/helper/must"
	"github.com/dosquad/mage/helper/paths"
	"github.com/fatih/color"
	"github.com/princjef/mageutil/shellcmd"
)

// ModTidy run go mod tidy on all workspaces.
func ModTidy(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Golang)

	var listout []byte
	{
		var err error
		listout, err = bins.Command("go list -m -f '{{.Dir}}'")
		must.PanicIfError(err, "unable to get go list output")
	}

	wd := paths.MustGetGitTopLevel()
	defer func() { _ = os.Chdir(wd) }()

	scanner := bufio.NewScanner(bytes.NewReader(listout))
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			//nolint:forbidigo // printing output
			fmt.Printf("%s %s\n", color.MagentaString(">"), color.New(color.Bold).Sprintf("cd %s", line))
			if err := os.Chdir(line); err != nil {
				return err
			}

			if err := shellcmd.Command("go mod tidy -go=" + build.GolangVersionRaw()).Run(); err != nil {
				return err
			}
		}
	}

	return nil
}
