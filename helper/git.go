package helper

import (
	"strings"

	"github.com/dosquad/go-giturl"
	"github.com/dosquad/mage/semver"
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

func GitURL() string {
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

	if u.Scheme == "git+ssh" {
		u.User = nil
		u.Scheme = "https"
	}

	u.Path = strings.TrimSuffix(u.Path, ".git")

	return u.String()
}

func GitHeadRev() string {
	return strings.TrimSpace(MustGetOutput(`git rev-parse --short HEAD`))
}

func GitHeadTagDescribe() string {
	out, err := shellcmd.Command("git describe --tags HEAD").Output()
	if err != nil {
		return "v0.0.0"
	}

	return strings.TrimSpace(string(out))
}

func GitSemver() string {
	out, err := shellcmd.Command("git describe --tags HEAD").Output()
	if err != nil {
		return "v0.0.0"
	}

	ver := strings.TrimSpace(string(out))
	if !strings.HasPrefix(ver, "v") {
		ver = "v" + ver
	}

	if idx := strings.Index(ver, "-"); idx >= 0 {
		ver = ver[:idx]
	}

	ver = semver.Canonical(ver)

	return ver
}

func GitSemverRaw() string {
	ver := GitSemver()
	if len(ver) > 1 {
		return ver[1:]
	}

	return ver
}

func SemverBumpPatch(v string) string {
	if semver.Canonical(v) != "" {
		return semver.MajorMinor(v) + "." + semver.IncDecimal(semver.Patch(v))
	}

	return ""
}
