package helper

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/na4ma4/go-permbits"
)

func fileExistsInPath(glob, rootDir string) bool {
	found := false
	_ = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if match, matchErr := filepath.Glob(filepath.Join(path, glob)); matchErr == nil {
				if len(match) > 0 {
					found = true
					return errors.New("file found")
				}
			}
		}

		return nil
	})

	return found
}

func FileExistsInPath(glob string, path ...string) bool {
	for _, p := range path {
		if fileExistsInPath(glob, p) {
			return true
		}
	}

	return false
}

func FileExists(path ...string) bool {
	for _, p := range path {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}

	return false
}

func FileWrite(data []byte, path string) error {
	var f *os.File
	{
		var err error
		f, err = os.Create(path)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func FileLineExists(filename, targetLine string) bool {
	var f *os.File
	{
		var err error
		f, err = os.Open(filename)
		if err != nil {
			return false
		}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == targetLine {
			return true
		}
	}

	return false
}

func FileLastRune(filename string) (rune, error) {
	var lastChar rune
	var f *os.File
	{
		var err error
		f, err = os.Open(filename)
		if err != nil {
			return lastChar, err
		}
	}

	rdr := bufio.NewReader(f)
	for {
		c, _, err := rdr.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return lastChar, err
		}

		lastChar = c
	}

	return lastChar, nil
}

func FileAppendLine(filename string, fileperm os.FileMode, line string) error {
	if fileperm == 0 {
		fileperm = permbits.MustString("a=r,u=rw")
	}
	var f *os.File
	{
		var err error
		f, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, fileperm)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	if _, err := f.WriteString(line + "\n"); err != nil {
		return err
	}

	return nil
}

func fileModTime(path string, defaultTime time.Time) time.Time {
	if st, err := os.Stat(path); err == nil {
		return st.ModTime()
	}

	return defaultTime
}

func TargetNeedRefresh(target string, src ...string) bool {
	targetAge := fileModTime(target, time.Time{})

	for _, item := range src {
		if fileModTime(item, time.Now()).After(targetAge) {
			return true
		}
	}

	return false
}

func FileNameModify(target string, fn func(string) []string) []string {
	return fn(target)
}

func FileNameModifyReplace(from string, to ...string) func(string) []string {
	return func(s string) []string {
		st := map[string]interface{}{}
		for _, item := range to {
			st[strings.Replace(s, from, item, 1)] = nil
		}

		out := []string{}
		for item := range st {
			out = append(out, item)
		}

		return out
	}
}
