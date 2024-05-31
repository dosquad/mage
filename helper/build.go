package helper

import (
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
