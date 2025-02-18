package build

import (
	"fmt"
	"strings"

	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/helper/must"
)

func KubernetesGetPodWithSelector(selector string) (string, error) {
	out, err := bins.Command(fmt.Sprintf(`kubectl get pod -l %s -o name`, selector))
	return strings.TrimSpace(string(out)), err
}

func KubernetesGetCurrentContext() (string, error) {
	out, err := bins.Command(`kubectl config view --minify -o jsonpath='{..namespace}'`)
	o := string(out)
	if o == "" {
		o = "default"
	}
	return o, must.IfErrorf("unable to execute 'kubectl config' command: %w", err)
}
