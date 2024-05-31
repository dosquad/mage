package helper

import (
	"fmt"
	"strings"

	"github.com/princjef/mageutil/shellcmd"
)

// func MustString(in string, err error) string {
// 	PanicIfError(err, "must not return error")
// 	return in
// }

func KubernetesGetPodWithSelector(selector string) (string, error) {
	out, err := shellcmd.Command(fmt.Sprintf(`kubectl get pod -l %s -o name`, selector)).Output()
	return strings.TrimSpace(string(out)), err
}

func KubernetesGetCurrentContext() (string, error) {
	out, err := shellcmd.Command(`kubectl config view --minify -o jsonpath='{..namespace}'`).Output()
	o := string(out)
	if o == "" {
		o = "default"
	}
	return o, err
}
