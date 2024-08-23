package helper

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/princjef/mageutil/shellcmd"
)

func CommandString(cmd string) (string, error) {
	out, err := Command(cmd)
	return strings.TrimSpace(string(out)), err
}

//nolint:forbidigo // print output.
func Command(cmd string) ([]byte, error) {
	fmt.Printf("%s %s\n", color.MagentaString(">"), color.New(color.Bold).Sprintf("%s", cmd))
	out, err := shellcmd.Command(cmd).Output()
	return out, err
}
