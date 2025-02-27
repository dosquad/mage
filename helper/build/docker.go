package build

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dosquad/mage/helper/ctxval"
	"github.com/dosquad/mage/helper/envs"
	"github.com/dosquad/mage/helper/paths"
	"gopkg.in/yaml.v3"
)

type DockerConfig struct {
	Image       string                 `yaml:"image,omitempty"`
	Platforms   []string               `yaml:"platforms,omitempty"`
	Tag         interface{}            `yaml:"tag,omitempty"`
	BlockedTags []string               `yaml:"blocked_tags,omitempty"`
	BuildArgs   map[string]string      `yaml:"build_args,omitempty"`
	Kubernetes  DockerConfigKubernetes `yaml:"kubernetes,omitempty"`
	Mirrord     DockerConfigMirrord    `yaml:"mirrord,omitempty"`
	Ignore      []string               `yaml:"ignore,omitempty"`
}

type DockerConfigMirrord struct {
	Targetless bool `yaml:"targetless,omitempty"`
}

type DockerConfigKubernetes struct {
	Deployment  string `yaml:"deployment,omitempty"`
	PodSelector string `yaml:"pod-selector,omitempty"`
}

type Platform struct {
	OS   string
	Arch string
}

func (d DockerConfig) GetTags() []string {
	if v := envs.GetEnv("DOCKER_TAGS", ""); v != "" {
		return strings.Split(v, " ")
	}
	if v := envs.GetEnv("DOCKER_TAG", ""); v != "" {
		return []string{v}
	}

	switch v := d.Tag.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	}

	if v := GitHeadTag(); v != "" {
		return []string{v}
	}

	if v := GitHeadTagDescribe(); v != "" {
		return []string{v}
	}

	return []string{"dev"}
}

func (d DockerConfig) GetImage() string {
	if v := envs.GetEnv("DOCKER_REPO", ""); v != "" {
		return v
	}

	if v := envs.GetEnv("DOCKER_IMAGE", ""); v != "" {
		return v
	}

	return d.Image
}

func (d DockerConfig) GetImageRef() string {
	tags := d.GetTags()
	image := d.GetImage()

	if len(tags) > 0 {
		return fmt.Sprintf("%s:%s", image, tags[0])
	}

	return image + ":dev"
}

func (d DockerConfig) IsBlocked(in string) bool {
	for _, tag := range d.BlockedTags {
		if strings.EqualFold(tag, in) {
			return true
		}
	}

	return false
}

func (d DockerConfig) OSArch() []Platform {
	out := []Platform{}
	for _, platform := range d.Platforms {
		sp := strings.SplitN(platform, "/", 2) //nolint:mnd // formatted string "os/arch".
		if len(sp) != 2 {                      //nolint:mnd // formatted string "os/arch".
			continue
		}

		out = append(out, Platform{
			OS:   sp[0],
			Arch: sp[1],
		})
	}

	return out
}

func (d DockerConfig) ArgsTag(tag string) (string, error) {
	if d.IsBlocked(tag) {
		return "", fmt.Errorf("tag is blocked: %s", tag)
	}

	if d.GetImage() == "" {
		return "", errors.New("image is not set")
	}

	return `--tag "` + d.GetImage() + `:` + tag + `"`, nil
}

func (d DockerConfig) Args(ctx context.Context) []string {
	out := []string{}

	if !ctxval.ContextDefaultValue[bool](ctx, ctxval.DockerLocalPlatform, false) {
		if len(d.Platforms) > 0 {
			out = append(out, "--platform "+strings.Join(d.Platforms, ","))
		}
	}

	if len(d.BuildArgs) > 0 {
		for key, arg := range d.BuildArgs {
			out = append(out, `--build-arg "`+key+`=`+arg+`"`)
		}
	}

	return out
}

func DockerLoadConfig() (*DockerConfig, error) {
	cfg := &DockerConfig{
		Platforms:   []string{"linux/amd64"},
		BlockedTags: []string{"dev"},
		BuildArgs: map[string]string{
			"VERSION": GitHeadTagDescribe(),
		},
		Kubernetes: DockerConfigKubernetes{},
		Ignore: []string{
			".makefiles",
			".github",
			".git",
		},
	}

	var f *os.File
	{
		var err error
		f, err = os.Open(paths.MustGetGitTopLevel(".docker.yml"))
		if err != nil {
			return cfg, err
		}
	}
	defer f.Close()

	{
		if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}
