package envs

import (
	"os"
	"strings"

	"github.com/princjef/mageutil/shellcmd"
)

// GetEnv returns an environment variable value if present, or a
// default value.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type CommandEnv struct {
	key string
	old string
}

func SetEnv(key, value string) (func() error, error) {
	e := &CommandEnv{
		key: key,
		old: os.Getenv(key),
	}

	if err := os.Setenv(key, value); err != nil {
		return nil, err
	}

	return e.Revert, nil
}

func (e *CommandEnv) Revert() error {
	return os.Setenv(e.key, e.old)
}

func GoEnv(key, fallback string) string {
	out, err := shellcmd.Command("go env " + key).Output()
	if err != nil {
		return fallback
	}

	return strings.TrimSpace(string(out))
}
