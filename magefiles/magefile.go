//go:build mage

package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/loga"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"

	//mage:import
	"github.com/dosquad/mage"
)

// TestLocal protoc, lint, test & build debug.
func TestLocal(ctx context.Context) {
	mg.CtxDeps(ctx, mage.Golang.Lint)

	mg.CtxDeps(ctx, DynamicDeps)
	mg.CtxDeps(ctx, mage.Golang.Test)
}

var Default = TestLocal

const (
	dyndepFileName = "artifacts/dyndep.test"
)

func DynamicDeps(ctx context.Context) error {
	_ = os.Remove(dyndepFileName)

	writeFile := func(ctx context.Context) error {
		loga.PrintInfo("Running Dynamic Dependency for dyndep.Test")
		_ = os.MkdirAll("artifacts", permbits.MustString("ug=rwx,o=rx"))
		_ = os.Remove(dyndepFileName)
		var f *os.File
		{
			var err error
			f, err = os.Create(dyndepFileName)
			if err != nil {
				return err
			}
		}
		defer f.Close()

		if _, err := io.WriteString(f, "this-is-a-test-file\n"); err != nil {
			return err
		}

		return nil
	}

	if paths.FileExists(dyndepFileName) {
		return fmt.Errorf("DynamicDeps: file exists when it should have been deleted: %s", dyndepFileName)
	}

	dyndep.Add(dyndep.Test, writeFile)
	mg.CtxDeps(ctx, mage.Golang.Test)

	if !paths.FileExists(dyndepFileName) {
		return fmt.Errorf("DynamicDeps: file does not exist when it should have been created: %s", dyndepFileName)
	}

	return nil
}
