package build

import (
	"errors"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dosquad/go-cliversion/makever"
	"github.com/dosquad/mage/helper/paths"
)

func boolToString(in bool) string {
	if in {
		return "true"
	}

	return "false"
}

// LDFlags returns the ldflags argument for `go build`.
func LDFlags(debug bool) []string {
	headTag := GitHeadTagDescribe()
	if headTag == "" {
		headTag = "0.0.0"
	}

	commonFlags := makever.LDFlags(
		makever.BuildDate(GitCommitTime().Format(time.RFC3339)),
		makever.BuildDebug(boolToString(debug)),
		makever.BuildMethod("magefiles"),
		makever.BuildGoVersion(runtime.Version()),
		makever.GitCommit(GitHash()),
		makever.GitRepo(GitURL()),
		makever.GitSlug(GitSlug()),
		makever.GitTag(GitSemver()),
		makever.GitExactTag(GitHeadTag()),
	)

	commonFlags = append(commonFlags,
		"-X main.commit="+GitHash(),
		"-X main.date="+time.Now().Format(time.RFC3339),
		"-X main.builtBy=magefiles",
		"-X main.repo="+GitURL(),
		"-X main.goVersion="+GolangVersionRaw(),
	)

	if debug {
		return append(commonFlags,
			"-X main.version="+headTag+"+debug",
			makever.BuildVersion(headTag+"+debug")(),
		)
	}

	return append(commonFlags,
		"-X main.version="+headTag,
		makever.BuildVersion(headTag)(),
		"-s",
		"-w",
	)
}

func FirstCommandName() (string, error) {
	path := paths.MustCommandPaths()
	if len(path) < 1 {
		return "", errors.New("command not found")
	}

	return filepath.Base(path[0]), nil
}
