package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/princjef/mageutil/shellcmd"
)

// Module is a module returned by `go list`.
type Module struct {
	Path       string       `json:"Path"`       // module path
	Query      string       `json:"Query"`      // version query corresponding to this version
	Version    string       `json:"Version"`    // module version
	Versions   []string     `json:"Versions"`   // available module versions
	Replace    *Module      `json:"Replace"`    // replaced by this module
	Time       *time.Time   `json:"Time"`       // time version was created
	Update     *Module      `json:"Update"`     // available update (with -u)
	Main       bool         `json:"Main"`       // is this the main module?
	Indirect   bool         `json:"Indirect"`   // module is only indirectly needed by main module
	Dir        string       `json:"Dir"`        // directory holding local copy of files, if any
	GoMod      string       `json:"GoMod"`      // path to go.mod file describing module, if any
	GoVersion  string       `json:"GoVersion"`  // go version used in module
	Retracted  []string     `json:"Retracted"`  // retraction information, if any (with -retracted or -u)
	Deprecated string       `json:"Deprecated"` // deprecation message, if any (with -u)
	Error      *ModuleError `json:"Error"`      // error loading module
	Origin     any          `json:"Origin"`     // provenance of module
	Reuse      bool         `json:"Reuse"`      // reuse of old module info is safe
}

// ModuleError is the possible error when loading a module.
type ModuleError struct {
	Err string `json:"Err"` // the error itself
}

// GolangListModules executes `go list -m all` and returned the parsed
// Module slice.
func GolangListModules() ([]Module, error) {
	out := []Module{}
	var modList []byte
	{
		var err error
		modList, err = shellcmd.Command(`go list -json -m all`).Output()
		if err != nil {
			return out, err
		}
	}

	d := json.NewDecoder(bytes.NewReader(modList))
	for d.More() {
		item := Module{}
		if err := d.Decode(&item); err != nil {
			return out, err
		}
		out = append(out, item)
	}

	return out, nil
}

// ProtobufTargetFunc returned by ProtobufTargets, each ProtobufTargetFunc
// returns a list of `*.proto` files in each directory.
type ProtobufTargetFunc func() []string

// ProtobufTargets returns a slice of functions that return a list of files
// in a specific directory
//
// example:
// root
// |-blah1 [test1.proto, test2.proto]
// \-blah2 [testa.proto, testb.proto]
//
// returns two functions,
// one that returns [/root/blah1/test1.proto, /root/blah1/test2.proto]
// and another that returns [/root/blah2/testa.proto /root/blah2/testb.proto].
func ProtobufTargets() []ProtobufTargetFunc {
	// Use a map to create a unique list of directories that contain `*.proto`.
	paths := map[string]interface{}{}
	for _, match := range FilesMatch(MustGetWD(), "*.proto") {
		paths[filepath.Dir(match)] = nil
	}

	// Create slice of functions with length of map.
	out := make([]ProtobufTargetFunc, len(paths))
	idx := 0
	for path := range paths {
		out[idx] = func() []string {
			matches, _ := filepath.Glob(filepath.Join(path, "*.proto"))
			return matches
		}
		idx++
	}

	return out
}

// ProtobufIncludePaths returns the `--proto_path` arguments for the protocol buffer
// gen command.
func ProtobufIncludePaths() []string {
	out := []string{}
	modules, _ := GolangListModules()
	for _, module := range modules {
		if len(FilesMatch(module.Dir, "*.proto")) > 0 {
			out = append(out,
				fmt.Sprintf("--proto_path=%s=%s", module.Path, module.Dir),
			)
		}
	}
	return out
}
