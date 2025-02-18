package helper

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/dosquad/mage/helper/bins"
	"github.com/dosquad/mage/helper/build"
	"github.com/dosquad/mage/helper/envs"
	"github.com/dosquad/mage/helper/must"
	"github.com/dosquad/mage/helper/paths"
)

type CommandTemplate struct {
	Debug bool

	CGO       string
	GoOS      string
	GoArch    string
	GoArm     string
	GoVersion string

	GitRev     string
	GitHash    string
	GitHeadTag string
	GitSlug    string

	LDFlags string

	CWD string

	OutputArtifact string
	TargetPath     string
	CommandDir     string
	CommandName    string
	HomeDir        string
}

func NewCommandTemplate(debug bool, commandDir string) *CommandTemplate {
	o := &CommandTemplate{
		Debug: debug,

		CGO:       envs.GetEnv("CGO_ENABLED", "0"),
		GoOS:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
		GoArm:     must.Must[string](bins.CommandString(`go env GOARM`)),
		GoVersion: runtime.Version(),

		GitRev:     build.GitHeadRev(),
		GitHash:    build.GitHash(),
		GitHeadTag: build.GitHeadTag(),
		GitSlug:    build.GitSlug(),

		LDFlags: strings.Join(build.LDFlags(debug), " "),

		CWD: paths.MustGetGitTopLevel(),
		// BaseDir: baseDir,

		CommandDir:  commandDir,
		CommandName: filepath.Base(commandDir),
		HomeDir:     paths.MustGetHomeDir(),
	}

	o.apply()

	return o
}

func (t *CommandTemplate) SetPlatform(os, arch, arm string) {
	t.GoOS = os
	t.GoArch = arch
	t.GoArm = arm

	t.apply()
}

func (t *CommandTemplate) apply() {
	if t.CommandName == "." {
		t.CommandName = "main"
	}

	releaseType := "release"

	if t.Debug {
		releaseType = "debug"
	}

	arch := t.GoArch
	if t.GoArm != "" {
		arch += "v" + t.GoArm
	}

	t.OutputArtifact = "artifacts/build/" + releaseType + "/" + t.GoOS + "/" + arch + "/" + t.CommandName

	if t.GitHeadTag == "" {
		t.GitHeadTag = "0.0.0"
	}

	if t.Debug {
		t.GitHeadTag += "+debug"
	}
}

func (t *CommandTemplate) Render(cmd string) (string, error) {
	t.apply()

	var tmpl *template.Template
	{
		var err error
		tmpl, err = template.New("").Parse(cmd)
		if err != nil {
			return "", fmt.Errorf("unable to parse command template: %w", err)
		}
	}

	sb := &strings.Builder{}
	if err := tmpl.Execute(sb, t); err != nil {
		return "", fmt.Errorf("unable to execute command template: %w", err)
	}

	return sb.String(), nil
}
