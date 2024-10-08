package mage

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper"
	"github.com/fatih/color"
	"github.com/princjef/mageutil/shellcmd"
)

// ModTidy run go mod tidy on all workspaces.
func ModTidy(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Golang)

	var listout []byte
	{
		var err error
		listout, err = helper.Command("go list -m -f '{{.Dir}}'")
		helper.PanicIfError(err, "unable to get go list output")
	}

	wd := helper.MustGetGitTopLevel()
	defer func() { _ = os.Chdir(wd) }()

	scanner := bufio.NewScanner(bytes.NewReader(listout))
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			//nolint:forbidigo // printing output
			fmt.Printf("%s %s\n", color.MagentaString(">"), color.New(color.Bold).Sprintf("cd %s", line))
			if err := os.Chdir(line); err != nil {
				return err
			}

			if err := shellcmd.Command("go mod tidy -go=" + helper.GolangVersionRaw()).Run(); err != nil {
				return err
			}
		}
	}

	return nil
}
