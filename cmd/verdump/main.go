package main

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"runtime/debug"

	"github.com/alecthomas/kong"
)

type Context struct {
	Debug   bool
	Version bool
}

type ModCmd struct {
	Path string `arg:"" name:"path" help:"path to binary" type:"path"`
}

//nolint:forbidigo // display version.
func (r *ModCmd) Run(_ *Context) error {
	out, err := exec.Command( //nolint:gosec // shellescape clean
		"go", "version", "-m",
		r.Path,
	).Output()
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`mod\s+(\S+)\s+((?:v|)[0-9.a-zA-Z()-]+)(?:\s+\S+|)`)
	matches := re.FindAllSubmatch(out, -1)

	if len(matches) > 0 && len(matches[0]) >= 3 {
		fmt.Printf("%s\n", matches[0][2])
		return nil
	}

	return errors.New("unable to parse output")
}

var cli struct {
	Debug   bool        `help:"Enable debug mode."`
	Version VersionFlag `name:"version" help:"Print version information and quit."`

	Mod ModCmd `cmd:"" help:"Get version from go module"`
}

type VersionFlag bool

func (v VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                       { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, _ kong.Vars) error {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		fmt.Println("version: " + buildInfo.Main.Version) //nolint:forbidigo // version display
		app.Exit(0)
		return nil
	}

	app.Exit(1)
	return errors.New("failed to retrieve module version")
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
