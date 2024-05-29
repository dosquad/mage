package helper

import (
	"fmt"

	"github.com/princjef/mageutil/shellcmd"
)

func MustGetOutput(cmd string) string {
	out, err := shellcmd.Command(cmd).Output()
	PanicIfError(err, fmt.Sprintf("unable to run command: %s", cmd))
	return string(out)
}
