package helper

import (
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

type CommandTemplate struct {
	Debug bool

	CGO    string
	GoOS   string
	GoArch string
	GoArm  string

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

		CGO:    GetEnv("CGO_ENABLED", "0"),
		GoOS:   runtime.GOOS,
		GoArch: runtime.GOARCH,
		GoArm:  strings.TrimSpace(MustGetOutput(`go env GOARM`)),

		GitRev:     GitHeadRev(),
		GitHash:    GitHash(),
		GitHeadTag: GitHeadTag(),
		GitSlug:    GitSlug(),

		LDFlags: strings.Join(LDFlags(debug), " "),

		CWD: MustGetWD(),
		// BaseDir: baseDir,

		CommandDir:  commandDir,
		CommandName: filepath.Base(commandDir),
		HomeDir:     MustGetHomeDir(),
	}

	if debug {
		o.OutputArtifact = "artifacts/build/debug/" + o.GoOS + "/" + o.GoArch + "/" + o.CommandName
	} else {
		o.OutputArtifact = "artifacts/build/release/" + o.GoOS + "/" + o.GoArch + "/" + o.CommandName
	}

	if o.GitHeadTag == "" {
		o.GitHeadTag = "0.0.0"
	}

	if debug {
		o.GitHeadTag += "+debug"
	}

	return o
}

func (t *CommandTemplate) Render(cmd string) (string, error) {
	var tmpl *template.Template
	{
		var err error
		tmpl, err = template.New("").Parse(cmd)
		if err != nil {
			return "", err
		}
	}

	sb := &strings.Builder{}
	if err := tmpl.Execute(sb, t); err != nil {
		return "", err
	}

	return sb.String(), nil
}
