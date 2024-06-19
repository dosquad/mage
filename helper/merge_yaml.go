package helper

import (
	"fmt"
)

func MergeYaml(leftFile, rightFile string) ([]byte, error) {
	out, err := Command(fmt.Sprintf(`yq -n 'load("%s") * load("%s")'`, leftFile, rightFile))
	return out, IfErrorf("unable to merge YAML files: %w", err)
}
