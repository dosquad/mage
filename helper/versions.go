package helper

import (
	"errors"
	"os"
	"strings"

	"github.com/na4ma4/go-permbits"
	"github.com/princjef/mageutil/shellcmd"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type VersionKey string

func (vk VersionKey) String() string {
	return string(vk)
}

func (vk VersionKey) Key() string {
	return cases.Upper(language.English).String(
		strings.ReplaceAll(vk.String(), "-", ""),
	) + "_VERSION"
}

const (
	GolangciLintVersion     VersionKey = "golangci-lint"
	GovulncheckVersion      VersionKey = "govulncheck"
	ProtocVersion           VersionKey = "protoc"
	ProtocGenGoVersion      VersionKey = "protoc-gen-go"
	ProtocGenGoGRPCVersion  VersionKey = "protoc-gen-go-grpc"
	ProtocGenGoTwirpVersion VersionKey = "protoc-gen-go-twirp"
	YQVersion               VersionKey = "yq"
)

const (
// golangciLintVersion     = "1.59.0"
// govulncheckVersion      = "latest"
// protocGenGoGRPCVersion  = "latest"
// protocGenGoTwirpVersion = "v8.1.3"
)

type VersionCache map[VersionKey]string

func MustVersionLoadCache() *VersionCache {
	cache, err := VersionLoadCache()
	if !errors.Is(err, os.ErrNotExist) {
		PanicIfError(err, "unable to load mage version config")
	}
	return cache
}

func VersionLoadCache() (*VersionCache, error) {
	cache := &VersionCache{
		GolangciLintVersion:     "",
		GovulncheckVersion:      "latest",
		ProtocVersion:           "",
		ProtocGenGoVersion:      "",
		ProtocGenGoGRPCVersion:  "latest",
		ProtocGenGoTwirpVersion: "",
		YQVersion:               "",
	}

	var f *os.File
	{
		var err error
		f, err = os.Open(MustGetArtifactPath(".versioncache.yaml"))
		if err != nil {
			return cache, err
		}
	}

	{
		if err := yaml.NewDecoder(f).Decode(&cache); err != nil {
			return cache, err
		}
	}

	return cache, nil
}

func (vc VersionCache) Save() error {
	{
		if !FileExists(MustGetArtifactPath()) {
			if err := os.MkdirAll(MustGetArtifactPath(), permbits.MustString("ug=rwx,o=rx")); err != nil {
				PrintWarning("unable to create artifact directory: %s", err)
				return err
			}
		}
	}
	var f *os.File
	{
		var err error
		f, err = os.Create(MustGetArtifactPath(".versioncache.yaml"))
		if err != nil {
			PrintWarning("unable to create version cache: %s", err)
			return err
		}
	}
	defer f.Close()

	if err := yaml.NewEncoder(f).Encode(vc); err != nil {
		PrintWarning("unable to encode version cache: %s", err)
		return err
	}

	return nil
}

func (vc VersionCache) SetVersion(key VersionKey, value string) string {
	vc[key] = value
	_ = vc.Save()
	return value
}

func (vc VersionCache) GetVersion(key VersionKey) string {
	if out, ok := vc[key]; ok && out != "" {
		return out
	}

	if v := GetEnv(key.String(), ""); v != "" {
		return vc.SetVersion(key, v)
	}

	switch key { //nolint:exhaustive // Don't have custom loaders for other keys yet.
	case ProtocVersion:
		return vc.SetVersion(key, vc.getProtocVersion())
	case ProtocGenGoVersion:
		return vc.SetVersion(key, vc.getProtobufVersion())
	case ProtocGenGoTwirpVersion:
		return vc.SetVersion(key, vc.getGithubVersion("twitchtv/twirp"))
	case GolangciLintVersion:
		return vc.SetVersion(key, vc.getGolangcilintVersion())
	case YQVersion:
		return vc.SetVersion(key, vc.getGithubVersion("mikefarah/yq"))
	}

	return ""
}

func (vc VersionCache) getGithubVersion(slug string) string {
	ver, _ := HTTPGetLatestGitHubVersion(slug)
	return ver
}

func (vc VersionCache) getGolangcilintVersion() string {
	if v := vc.getGithubVersion("golangci/golangci-lint"); v != "" {
		return strings.TrimPrefix(v, "v")
	}

	return ""
}

func (vc VersionCache) getProtocVersion() string {
	protocVer, err := HTTPGetLatestGitHubVersion("protocolbuffers/protobuf")
	if err != nil {
		PrintWarning("Protocol Buffer Error: %s", err)
		return "latest"
	}

	return strings.TrimPrefix(protocVer, "v")
}

func (vc VersionCache) getProtobufVersion() string {
	ver, err := shellcmd.Command(`go list -f '{{.Version}}' -m "google.golang.org/protobuf"`).Output()
	if err != nil {
		PrintWarning("Warning: did not find google.golang.org/protobuf in go.mod, defaulting to latest")
		return "latest"
	}

	return strings.TrimSpace(string(ver))
}
