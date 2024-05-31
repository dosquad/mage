package helper

import (
	"errors"
	"path/filepath"
	"time"
)

func LDFlags(debug bool) []string {
	headTag := GitHeadTagDescribe()
	if headTag == "" {
		headTag = "0.0.0"
	}

	commonFlags := []string{
		"-X main.commit=" + GitHash(),
		"-X main.date=" + time.Now().Format(time.RFC3339),
		"-X main.builtBy=magefiles",
		"-X main.repo=" + GitURL(),
		"-X main.goVersion=" + GolangVersionRaw(),
	}

	if debug {
		return append(commonFlags,
			"-X main.version="+headTag+"+debug",
		)
	}

	return append(commonFlags,
		"-X main.version="+headTag,
		"-s",
		"-w",
	)
}

func MustFirstCommandName() string {
	cmd, err := FirstCommandName()
	PanicIfError(err, "unable to find command name")
	return cmd
}

func FirstCommandName() (string, error) {
	paths := MustCommandPaths()
	if len(paths) < 1 {
		return "", errors.New("command not found")
	}

	return filepath.Base(paths[0]), nil
}
