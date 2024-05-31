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

func ArgsFromAny(in []any) []string {
	out := []string{}
	for _, item := range in {
		switch v := item.(type) {
		case string:
			out = append(out, v)
		case []string:
			out = append(out, v...)
		}
	}

	return out
}
