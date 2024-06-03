package helper

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/na4ma4/go-permbits"
)

func GetRelativePath(path string) (string, bool) {
	root := filepath.Clean(MustGetWD())
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
func MustGetGoBin() string {
	if goBin := GetEnv("GOBIN", ""); goBin != "" {
		return goBin
	}

	if goPath := GetEnv("GOPATH", ""); goPath != "" {
		return filepath.Join(goPath, "bin")
	}

	return MustGetArtifactPath("bin")
}

func MustGetProtobufPath() string {
	protobufPath := MustGetArtifactPath("protobuf", "bin")
	// _ = os.MkdirAll(protobufPath, permbits.MustString("ug=rwx,o=rx"))
	return protobufPath
}

// MustGetArtifactPath get artifact directory or panic if unable to.
func MustGetArtifactPath(path ...string) string {
	return MustGetWD(append([]string{"artifacts"}, path...)...)
}

// MustGetWD get working directory or panic if unable to.
func MustGetWD(path ...string) string {
	wd := Must[string](os.Getwd())
	return filepath.Join(append([]string{wd}, path...)...)
}

// MustGetHomeDir get user home directory or panic if unable to.
func MustGetHomeDir(path ...string) string {
	homeDir := Must[string](os.UserHomeDir())
	return filepath.Join(append([]string{homeDir}, path...)...)
}

// MustMakeDir make directory, including parents if required.
//
// if unable to make directory then panic.
func MustMakeDir(path string, fileperm os.FileMode) {
	if fileperm == 0 {
		fileperm = os.ModePerm
	}
	PanicIfError(os.MkdirAll(path, fileperm), fmt.Sprintf("unable to make dir: [%s]", path))
}

// // mustGetCoverageOutPath get output directory for generated coverage report.
// //
// // if the directory does not exist then it will be created.
// //
// // if the output directory cannot be created then panic.
// func mustGetCoverageOutPath() string {
// 	outPath := filepath.Join(mustGetWD(), "artifacts", "coverage")

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

func MustCommandPaths() []string {
	if FileExists(MustGetWD("cmd")) {
		out := []string{}
		wd := MustGetWD("cmd")
		_ = filepath.WalkDir(wd, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if !d.IsDir() {
				return nil
			}

			if filepath.Dir(path) != MustGetWD("cmd") {
				return nil
			}

			out = append(out,
				fmt.Sprintf("./cmd/%s", d.Name()),
			)

			return nil
		})

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
	vscPath := filepath.Join(append([]string{MustGetWD()}, path...)...)

	MustMakeDir(vscPath, permbits.MustString("ug=rwx,o=rx"))

	return vscPath
}
