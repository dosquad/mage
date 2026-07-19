package mage

import (
	"archive/zip"
	"context"
	"io"
	"os"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper"
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/loga"

	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
)

// Lambda creates lambda artifacts.
func (Build) Lambda(ctx context.Context) error {
	// Lambda builds only support the following platforms.
	_ = os.Setenv("PLATFORMS", "linux/arm64,linux/amd64")

	dyndep.CtxDeps(ctx, dyndep.Build)
	mg.CtxDeps(ctx, Build.Release)

	pathList := paths.MustCommandPaths()

	paths.MustMakeDir(paths.MustGetArtifactPath("lambda"), permbits.MustString("ug=rwx,o=rx"))

	for _, cmdPath := range pathList {
		ct := helper.NewCommandTemplate(false, cmdPath)
		ct.AdditionalArtifacts = dyndep.GetArchiveSideloadDeps(ctx)
		loga.PrintDebugf("ct.AdditionalArtifacts: %+v", ct.AdditionalArtifacts)

		if err := buildPlatformIterator(ctx, ct, func(_ context.Context, ct *helper.CommandTemplate) error {
			// buildArtifact(ctx, ct)
			loga.PrintInfof("Create Archive for %s/%s (%s)",
				ct.GoOS, ct.GoArch,
				ct.OutputArtifact,
			)

			if err := buildLambdaZip(ctx, ct); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func buildLambdaZip(ctx context.Context, ct *helper.CommandTemplate) error {
	zipPath := paths.MustGetArtifactPath("lambda", ct.CommandName+"_"+ct.GoOS+"_"+ct.GoArch+"-lambda-bootstrap.zip")
	loga.PrintDebugf("buildLambdaZip(%s, %s, %s): %s", ct.CommandName, ct.GoOS, ct.GoArch, zipPath)

	var srcFile *os.File
	{
		var err error
		srcFile, err = os.Open(ct.OutputArtifact)
		if err != nil {
			return err
		}
		defer srcFile.Close()
	}

	var f *os.File
	{
		var err error
		f, err = os.Create(zipPath)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	zipf := zip.NewWriter(f)
	defer zipf.Close()

	var bsFile io.Writer
	{
		var err error
		bsFile, err = zipf.CreateHeader(lambdaZipHeader(ctx))
		if err != nil {
			return err
		}
	}

	if _, err := io.Copy(bsFile, srcFile); err != nil {
		return err
	}

	if err := zipf.Close(); err != nil {
		return err
	}

	return nil
}

func lambdaZipHeader(_ context.Context) *zip.FileHeader {
	h := &zip.FileHeader{
		Name: "bootstrap",
	}

	// Define unix permissions for the file
	mode := permbits.MustString("a=rx,u+w")

	h.ExternalAttrs = uint32(mode) << 16 //nolint:mnd // Store permissions in the top 16 bits

	//nolint:mnd // Set CreatorVersion to signal Unix attributes (3 << 8 represents Unix OS)
	h.CreatorVersion = 3 << 8

	return h
}
