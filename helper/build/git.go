package build

import (
	"strconv"
	"strings"
	"time"

	"github.com/dosquad/go-giturl"
	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/helper/must"
	"github.com/dosquad/mage/semver"
)

func GitHash() string {
	out, err := bins.Command("git show -s --format=%h")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func GitCommitTime() time.Time {
	var out string
	{
		var err error
		out, err = bins.CommandString(`git log -1 --format="%at"`)
		if err != nil {
			return time.Time{}
		}
	}

	if v, err := strconv.ParseInt(out, 10, 64); err == nil {
		return time.Unix(v, 0)
	}

	return time.Time{}
}

func GitHeadTag() string {
	out, err := bins.Command("git describe --tags --exact-match HEAD 2>/dev/null")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func GitSlug() string {
	out, err := bins.Command("git config --get remote.origin.url")
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
	out, err := bins.Command("git config --get remote.origin.url")
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
	return must.Must[string](bins.CommandString(`git rev-parse --short HEAD`))
}

func GitHeadTagDescribe() string {
	out, err := bins.Command("git describe --tags HEAD")
	if err != nil {
		return "v0.0.0"
	}

	return strings.TrimSpace(string(out))
}

func GitSemver() string {
	out, err := bins.Command("git describe --tags HEAD")
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
