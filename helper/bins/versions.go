package bins

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/dosquad/mage/helper/envs"
	"github.com/dosquad/mage/helper/must"
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/helper/web"
	"github.com/dosquad/mage/loga"
	"github.com/na4ma4/go-permbits"
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
	latestTag = "latest"
)

const (
	GolangciLintVersion          VersionKey = "golangci-lint"
	GovulncheckVersion           VersionKey = "govulncheck"
	ProtocVersion                VersionKey = "protoc"
	ProtocGenGoVersion           VersionKey = "protoc-gen-go"
	ProtocGenGoGRPCVersion       VersionKey = "protoc-gen-go-grpc"
	ProtocGenGoConnectVersion    VersionKey = "protoc-gen-connect-go"
	ProtocGenGoTwirpVersion      VersionKey = "protoc-gen-go-twirp"
	YQVersion                    VersionKey = "yq"
	BufVersion                   VersionKey = "buf"
	GoreleaserVersion            VersionKey = "goreleaser"
	WireVersion                  VersionKey = "wire"
	VerdumpVersion               VersionKey = "verdump"
	KubeControllerGenVersion     VersionKey = "kubernetes-controller-gen"
	KustomizeVersion             VersionKey = "kustomize"
	KubeControllerEnvTestVersion VersionKey = "kubernetes-controller-env-test"
	CFSSLVersion                 VersionKey = "cfssl"
	VGTVersion                   VersionKey = "vgt"
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
		must.PanicIfError(err, "unable to load mage version config")
	}
	return cache
}

func VersionLoadCache() (*VersionCache, error) {
	cache := &VersionCache{
		BufVersion:                   "",
		GolangciLintVersion:          "",
		GovulncheckVersion:           latestTag,
		ProtocVersion:                "",
		ProtocGenGoVersion:           "",
		ProtocGenGoGRPCVersion:       "",
		ProtocGenGoTwirpVersion:      "",
		ProtocGenGoConnectVersion:    "",
		YQVersion:                    "",
		GoreleaserVersion:            "",
		WireVersion:                  "",
		VerdumpVersion:               "",
		KustomizeVersion:             "",
		KubeControllerGenVersion:     "",
		KubeControllerEnvTestVersion: "",
		CFSSLVersion:                 "",
		VGTVersion:                   "",
	}

	var f *os.File
	{
		var err error
		f, err = os.Open(paths.MustGetArtifactPath(".versioncache.yaml"))
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
		if !paths.FileExists(paths.MustGetArtifactPath()) {
			if err := os.MkdirAll(paths.MustGetArtifactPath(), permbits.MustString("ug=rwx,o=rx")); err != nil {
				loga.PrintWarning("unable to create artifact directory: %s", err)
				return err
			}
		}
	}
	var f *os.File
	{
		var err error
		f, err = os.Create(paths.MustGetArtifactPath(".versioncache.yaml"))
		if err != nil {
			loga.PrintWarning("unable to create version cache: %s", err)
			return err
		}
	}
	defer f.Close()

	if err := yaml.NewEncoder(f).Encode(vc); err != nil {
		loga.PrintWarning("unable to encode version cache: %s", err)
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
	if v, ok := vc[key]; ok && v != "" {
		return v
	}

	if v := envs.GetEnv(key.String(), ""); v != "" {
		return vc.SetVersion(key, v)
	}

	switch key { //nolint:exhaustive // Don't have custom loaders for other keys yet.
	case ProtocVersion:
		return vc.SetVersion(key, vc.getProtocVersion())
	case ProtocGenGoVersion:
		return vc.SetVersion(key, vc.getProtobufVersion())
	case ProtocGenGoTwirpVersion:
		return vc.SetVersion(key, vc.getGithubVersion("twitchtv/twirp"))
	case ProtocGenGoConnectVersion:
		return vc.SetVersion(key, vc.getGithubVersion("connectrpc/connect-go"))
	case ProtocGenGoGRPCVersion:
		v, err := web.HTTPGetLatestGitHubReleaseMatchingTag("grpc/grpc-go", regexp.MustCompile(`^cmd/protoc-gen-go-grpc/`))
		if err != nil {
			return latestTag
		}

		return vc.SetVersion(key, strings.TrimPrefix(v, "cmd/protoc-gen-go-grpc/"))
	case GolangciLintVersion:
		return vc.SetVersion(key, vc.getGolangcilintVersion())
	case YQVersion:
		return vc.SetVersion(key, vc.getGithubVersion("mikefarah/yq"))
	case BufVersion:
		return vc.SetVersion(key, vc.getGithubVersion("bufbuild/buf"))
	case GoreleaserVersion:
		ver := vc.SetVersion(key, strings.TrimPrefix(vc.getGithubVersion("goreleaser/goreleaser"), "v"))
		return ver
	case WireVersion:
		return vc.SetVersion(key, vc.getGithubVersion("google/wire"))
	case VerdumpVersion:
		return vc.SetVersion(key, vc.getGithubVersion("dosquad/mage"))
	case KubeControllerGenVersion:
		return vc.SetVersion(key, vc.getGithubVersion("kubernetes-sigs/controller-tools"))
	case KustomizeVersion:
		if after, found := strings.CutPrefix(
			vc.getGithubVersion("kubernetes-sigs/kustomize"),
			"kustomize/v",
		); found {
			return vc.SetVersion(key, after)
		}
	case KubeControllerEnvTestVersion:
		return vc.SetVersion(key, vc.getGithubVersion("kubernetes-sigs/controller-runtime"))
	case CFSSLVersion:
		return vc.SetVersion(key, vc.getGithubVersion("cloudflare/cfssl"))
	case VGTVersion:
		return vc.SetVersion(key, vc.getGithubVersion("roblaszczak/vgt"))
	}

	return ""
}

func (vc VersionCache) getGithubVersion(slug string) string {
	ver, _ := web.HTTPGetLatestGitHubVersion(slug)
	return ver
}

func (vc VersionCache) getGolangcilintVersion() string {
	if v := vc.getGithubVersion("golangci/golangci-lint"); v != "" {
		return strings.TrimPrefix(v, "v")
	}

	return ""
}

func (vc VersionCache) getProtocVersion() string {
	protocVer, err := web.HTTPGetLatestGitHubVersion("protocolbuffers/protobuf")
	if err != nil {
		loga.PrintWarning("Protocol Buffer Error: %s", err)
		return "latest"
	}

	return strings.TrimPrefix(protocVer, "v")
}

func (vc VersionCache) getProtobufVersion() string {
	ver, err := Command(`go list -f '{{.Version}}' -m "google.golang.org/protobuf"`)
	if err != nil {
		loga.PrintWarning("Warning: did not find google.golang.org/protobuf in go.mod, defaulting to latest")
		return "latest"
	}

	return strings.TrimSpace(string(ver))
}
