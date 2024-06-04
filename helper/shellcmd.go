package helper

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/princjef/mageutil/shellcmd"
)

func CommandString(cmd string) (string, error) {
	out, err := Command(cmd)
	return strings.TrimSpace(string(out)), err
}

func Command(cmd string) ([]byte, error) {
	if mg.Debug() || mg.Verbose() {
		fmt.Printf("%s %s\n", color.MagentaString(">"), color.New(color.Bold).Sprintf(cmd))
	}
	out, err := shellcmd.Command(cmd).Output()
	return out, err
}
