package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/princjef/mageutil/shellcmd"
)

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

type ModuleError struct {
	Err string `json:"Err"` // the error itself
}

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

func GetProtobufVersion() string {
	ver, err := shellcmd.Command(`go list -f '{{.Version}}' -m "google.golang.org/protobuf"`).Output()
	if err != nil {
		PrintWarning("Warning: did not find google.golang.org/protobuf in go.mod, defaulting to latest")
		return "latest"
	}

	return strings.TrimSpace(string(ver))
}

func ProtobufTargets() []string {
	return FilesMatch(MustGetWD(), "*.proto")
}

// func ProtobufTargets() []string {
// 	paths := map[string]interface{}{}
// 	for _, match := range FilesMatch(MustGetWD(), "*.proto") {
// 		paths[filepath.Dir(match)] = nil
// 	}

// 	out := make([]string, len(paths))
// 	idx := 0
// 	for path := range paths {
// 		out[idx] = filepath.Join(path, "*.proto")
// 		idx += 1
// 	}

// 	return out
// }

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

// func ProtobufFileNameModify(in string) string {
// 	return strings.Replace(in, ".proto", ".pb.go", 1)
// }
