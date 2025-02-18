package helper

import (
	"fmt"

	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/helper/must"
)

func MergeYaml(leftFile, rightFile string) ([]byte, error) {
	out, err := bins.Command(fmt.Sprintf(`yq -n 'load("%s") * load("%s")'`, leftFile, rightFile))
	return out, must.IfErrorf("unable to merge YAML files: %w", err)
}
