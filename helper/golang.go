package helper

import (
	"strings"
)

func GetModuleName() (string, error) {
	module, err := CommandString(`go list -m`)
	if err != nil {
		return "", err
	}

	sp := strings.SplitN(module, "\n", 2) //nolint:mnd // `|head -n1` equivalent

	return sp[0], nil
}

func GolangVersion() string {
	return Must[string](CommandString("go env GOVERSION"))
}

func GolangVersionRaw() string {
	return strings.TrimPrefix(GolangVersion(), "go")
}
