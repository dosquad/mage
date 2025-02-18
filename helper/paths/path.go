package paths

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dosquad/mage/helper/envs"
	"github.com/dosquad/mage/helper/must"
	"github.com/na4ma4/go-permbits"
	"github.com/princjef/mageutil/shellcmd"
)

func GetRelativePath(path string) (string, bool) {
	root := filepath.Clean(MustGetGitTopLevel())
	path = filepath.Clean(path)

	after, ok := strings.CutPrefix(path, root)
	if ok {
		return "." + after, true
	}

	return after, false
}

// MustGetGoBin get GOBIN.
//
// if GOBIN is not set then attempt to derive its
// value using GOPATH.
//
// if GOPATH is not set, use (pwd)/artifacts/bin .
func MustGetGoBin(path ...string) string {
	if goBin := envs.GetEnv("GOBIN", ""); goBin != "" {
		return filepath.Join(append([]string{goBin}, path...)...)
	}

	if goPath := envs.GetEnv("GOPATH", ""); goPath != "" {
		return filepath.Join(append([]string{goPath, "bin"}, path...)...)
	}

	return MustGetArtifactPath(append([]string{"bin"}, path...)...)
}

func MustGetProtobufPath() string {
	protobufPath := MustGetArtifactPath("protobuf", "bin")
	// _ = os.MkdirAll(protobufPath, permbits.MustString("ug=rwx,o=rx"))
	return protobufPath
}

// MustGetArtifactPath get artifact directory or panic if unable to.
func MustGetArtifactPath(path ...string) string {
	return MustGetGitTopLevel(append([]string{"artifacts"}, path...)...)
}

// MustGetWD get working directory or panic if unable to.
func MustGetWD(path ...string) string {
	wd := must.Must[string](os.Getwd())
	return filepath.Join(append([]string{wd}, path...)...)
}

// MustGetHomeDir get user home directory or panic if unable to.
func MustGetHomeDir(path ...string) string {
	homeDir := must.Must[string](os.UserHomeDir())
	return filepath.Join(append([]string{homeDir}, path...)...)
}

// MustMakeDir make directory, including parents if required.
//
// if unable to make directory then panic.
func MustMakeDir(path string, fileperm os.FileMode) {
	if fileperm == 0 {
		fileperm = os.ModePerm
	}
	must.PanicIfError(os.MkdirAll(path, fileperm), fmt.Sprintf("unable to make dir: [%s]", path))
}

// // mustGetCoverageOutPath get output directory for generated coverage report.
// //
// // if the directory does not exist then it will be created.
// //
// // if the output directory cannot be created then panic.
// func mustGetCoverageOutPath() string {
// 	outPath := filepath.Join(MustGetGitTopLevel(), "artifacts", "coverage")

// 	mustMakeDir(outPath)

// 	return outPath
// }

func FilesMatch(baseDir, pattern string) []string {
	matches := []string{}
	_ = filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			switch d.Name() {
			case "artifacts", "protobuf":
				return filepath.SkipDir
			}

			return nil
		}

		if match, err := filepath.Match(pattern, d.Name()); err == nil && match {
			matches = append(matches, path)
		}

		return nil
	})

	return matches
}

//nolint:nestif // handling ./cmd/main.go as well as ./cmd/foo/main.go
func MustCommandPaths() []string {
	if FileExists(MustGetGitTopLevel("cmd")) {
		out := []string{}
		wd := MustGetGitTopLevel("cmd")
		fallback := []string{}
		_ = filepath.WalkDir(wd, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if !d.IsDir() {
				localPath, err := filepath.Rel(wd, path)
				if err != nil {
					return nil //nolint:nilerr // continue processing
				}
				localPath = "./cmd/" + filepath.Dir(localPath)
				fallback = append(fallback, localPath)
				return nil
			}

			if filepath.Dir(path) != MustGetGitTopLevel("cmd") {
				return nil
			}

			out = append(out, "./cmd/"+d.Name())

			return nil
		})

		if len(out) == 0 {
			if len(fallback) > 0 {
				return []string{
					fallback[0],
				}
			}
		}

		return out
	}

	return []string{}
}

// MustGetVSCodePath get generated .vscode directory.
//
// if the directory does not exist then it will be created.
//
// if the directory cannot be created then panic.
func MustGetVSCodePath(path ...string) string {
	vscPath := MustGetGitTopLevel(append([]string{".vscode"}, path...)...)

	MustMakeDir(vscPath, permbits.MustString("ug=rwx,o=rx"))

	return vscPath
}

func MustGetGitTopLevel(path ...string) string {
	return must.Must[string](GetGitTopLevel(path...))
}

//nolint:gochecknoglobals // caching output from git command.
var gitTopLevel *string

func GetGitTopLevel(path ...string) (string, error) {
	if gitTopLevel != nil {
		return filepath.Join(append([]string{*gitTopLevel}, path...)...), nil
	}

	out, err := shellcmd.Command(`git rev-parse --show-toplevel`).Output()
	if err != nil {
		return "", err
	}

	localPath := string(out)
	localPath = strings.TrimSpace(localPath)

	gitTopLevel = &localPath

	return filepath.Join(append([]string{localPath}, path...)...), nil
}
