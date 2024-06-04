package helper

// func GetOutput(cmd string) (string, error) {
// 	out, err := shellcmd.Command(cmd).Output()
// 	return strings.TrimSpace(string(out)), err
// }

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
