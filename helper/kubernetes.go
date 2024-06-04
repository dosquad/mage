package helper

import (
	"fmt"
	"strings"
)

func KubernetesGetPodWithSelector(selector string) (string, error) {
	out, err := Command(fmt.Sprintf(`kubectl get pod -l %s -o name`, selector))
	return strings.TrimSpace(string(out)), err
}

func KubernetesGetCurrentContext() (string, error) {
	out, err := Command(`kubectl config view --minify -o jsonpath='{..namespace}'`)
	o := string(out)
	if o == "" {
		o = "default"
	}
	return o, IfErrorf("unable to execute 'kubectl config' command: %w", err)
}
