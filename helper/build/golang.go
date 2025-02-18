package build

import (
	"strings"

	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/helper/must"
)

func GetModuleName() (string, error) {
	module, err := bins.CommandString(`go list -m`)
	if err != nil {
		return "", err
	}

	sp := strings.SplitN(module, "\n", 2) //nolint:mnd // `|head -n1` equivalent

	return sp[0], nil
}

func GolangVersion() string {
	return must.Must[string](bins.CommandString("go env GOVERSION"))
}

func GolangVersionRaw() string {
	return strings.TrimPrefix(GolangVersion(), "go")
}
