package mage

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/dosquad/mage/helper"
	"github.com/fatih/color"
	"github.com/princjef/mageutil/shellcmd"
)

// ModTidy run go mod tidy on all workspaces.
func ModTidy() error {
	var listout []byte
	{
		var err error
		listout, err = shellcmd.Command("go list -m -f '{{.Dir}}'").Output()
		helper.PanicIfError(err, "unable to get go list output")
	}

	wd := helper.MustGetWD()
	defer func() { _ = os.Chdir(wd) }()

	scanner := bufio.NewScanner(bytes.NewReader(listout))
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			//nolint:forbidigo // printing output
			fmt.Printf("%s %s\n", color.MagentaString(">"), color.New(color.Bold).Sprintf("cd "+line))
			if err := os.Chdir(line); err != nil {
				return err
			}

			if err := shellcmd.Command("go mod tidy").Run(); err != nil {
				return err
			}
		}
	}

	return nil
}
