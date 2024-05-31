package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type DockerConfig struct {
	Image       string            `yaml:"image,omitempty"`
	Platforms   []string          `yaml:"platforms,omitempty"`
	Tag         interface{}       `yaml:"tag,omitempty"`
	BlockedTags []string          `yaml:"blocked_tags,omitempty"`
	BuildArgs   map[string]string `yaml:"build_args,omitempty"`
}

type BuildPlatform struct {
	OS   string
	Arch string
}

func (d DockerConfig) GetTags() []string {
	if v := GetEnv("DOCKER_TAGS", ""); v != "" {
		return strings.Split(v, " ")
	}
	if v := GetEnv("DOCKER_TAG", ""); v != "" {
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
	if v := GetEnv("DOCKER_REPO", ""); v != "" {
		return v
	}

	if v := GetEnv("DOCKER_IMAGE", ""); v != "" {
		return v
	}

	return d.Image
}

func (d DockerConfig) IsBlocked(in string) bool {
	for _, tag := range d.BlockedTags {
		if strings.EqualFold(tag, in) {
			return true
		}
	}

	return false
}

func (d DockerConfig) OSArch() []BuildPlatform {
	out := []BuildPlatform{}
	for _, platform := range d.Platforms {
		sp := strings.SplitN(platform, "/", 2) //nolint:mnd // formatted string "os/arch".
		if len(sp) != 2 {                      //nolint:mnd // formatted string "os/arch".
			continue
		}

		out = append(out, BuildPlatform{
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

func (d DockerConfig) Args() []string {
	out := []string{}

	if len(d.Platforms) > 0 {
		out = append(out, "--platform "+strings.Join(d.Platforms, ","))
	}

	if len(d.BuildArgs) > 0 {
		for key, arg := range d.BuildArgs {
			out = append(out, `--build-arg "`+key+`=`+arg+`"`)
		}
	}

	return out
}

func MustDockerLoadConfig() *DockerConfig {
	cfg, err := DockerLoadConfig()
	if !errors.Is(err, os.ErrNotExist) {
		PanicIfError(err, "unable to load mage docker config")
	}
	return cfg
}

func DockerLoadConfig() (*DockerConfig, error) {
	cfg := &DockerConfig{
		Platforms:   []string{"linux/amd64"},
		BlockedTags: []string{"dev"},
		BuildArgs: map[string]string{
			"VERSION": GitHeadTagDescribe(),
		},
	}

	var f *os.File
	{
		var err error
		f, err = os.Open(MustGetWD(".docker.yml"))
		if err != nil {
			return cfg, err
		}
	}

	{
		if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}
