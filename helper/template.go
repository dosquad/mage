package helper

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
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

		CGO:       GetEnv("CGO_ENABLED", "0"),
		GoOS:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
		GoArm:     Must[string](CommandString(`go env GOARM`)),
		GoVersion: runtime.Version(),

		GitRev:     GitHeadRev(),
		GitHash:    GitHash(),
		GitHeadTag: GitHeadTag(),
		GitSlug:    GitSlug(),

		LDFlags: strings.Join(LDFlags(debug), " "),

		CWD: MustGetGitTopLevel(),
		// BaseDir: baseDir,

		CommandDir:  commandDir,
		CommandName: filepath.Base(commandDir),
		HomeDir:     MustGetHomeDir(),
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

	if t.Debug {
		t.OutputArtifact = "artifacts/build/debug/" + t.GoOS + "/" + t.GoArch + "/" + t.CommandName
	} else {
		t.OutputArtifact = "artifacts/build/release/" + t.GoOS + "/" + t.GoArch + "/" + t.CommandName
	}

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
