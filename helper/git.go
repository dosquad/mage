package helper

import (
	"strings"

	"github.com/dosquad/go-giturl"
	"github.com/princjef/mageutil/shellcmd"
)

func GitHash() string {
	out, err := shellcmd.Command("git show -s --format=%h").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func GitHeadTag() string {
	out, err := shellcmd.Command("git describe --tags --exact-match HEAD 2>/dev/null").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func GitSlug() string {
	out, err := shellcmd.Command("git config --get remote.origin.url").Output()
	if err != nil {
		return ""
	}

	origin := strings.TrimSpace(string(out))
	u, err := giturl.Parse(origin)
	if err != nil {
		origin = strings.TrimSuffix(origin, ".git")
		return origin
	}

	return u.Slug()
}

func GitHeadRev() string {
	return strings.TrimSpace(MustGetOutput(`git rev-parse --short HEAD`))
}
