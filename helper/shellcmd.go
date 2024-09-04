package helper

import (
	"strings"

	"github.com/dosquad/mage/loga"
	"github.com/princjef/mageutil/shellcmd"
)

func CommandString(cmd string) (string, error) {
	out, err := Command(cmd)
	return strings.TrimSpace(string(out)), err
}

//nolint:forbidigo // print output.
func Command(cmd string) ([]byte, error) {
	loga.PrintCommand(cmd)
	out, err := shellcmd.Command(cmd).Output()
	return out, err
}
