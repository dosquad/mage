package helper

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/princjef/mageutil/shellcmd"
)

//nolint:forbidigo // printing output
func MergeYaml(leftFile, rightFile string) ([]byte, error) {
	cmd := fmt.Sprintf(`yq -n 'load("%s") * load("%s")'`, leftFile, rightFile)
	fmt.Printf("%s %s\n", color.MagentaString(">"), color.New(color.Bold).Sprintf(cmd))
	return shellcmd.Command(cmd).Output()
}
