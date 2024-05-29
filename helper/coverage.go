package helper

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/na4ma4/go-permbits"
)

// FilterCoverageOutput filters the coverage output by removing
// any lines that contain generated protobuf files `.pg.go`.
func FilterCoverageOutput(filename string) error {
	var f *os.File
	{
		var err error
		f, err = os.Open(filename)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), ".pb.go") {
			if _, err := buf.Write(scanner.Bytes()); err != nil {
				return err
			}

			if _, err := buf.WriteString("\n"); err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := os.WriteFile(filename, buf.Bytes(), permbits.MustString("a=rw")); err != nil {
		return err
	}

	return nil
}
