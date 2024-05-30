package helper

import (
	"strings"

	"github.com/princjef/mageutil/shellcmd"
)

func GetModuleName() (string, error) {
	module, err := shellcmd.Command(`go list -m`).Output()
	if err != nil {
		return "", err
	}

	sp := strings.SplitN(string(module), "\n", 2) //nolint:mnd // `|head -n1` equivalent

	return sp[0], nil
}

func MustModuleName() string {
	module, err := GetModuleName()
	PanicIfError(err, "unable to retrieve module name")
	return module
}
