package mage

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dosquad/mage/dyndep"
	"github.com/dosquad/mage/helper"
	"github.com/dosquad/mage/helper/build"
	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/loga"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
)

// Archive creates the archive artifacts.
func (Build) Archive(ctx context.Context) error {
	dyndep.CtxDeps(ctx, dyndep.Build)
	mg.CtxDeps(ctx, Build.Release)

	pathList := paths.MustCommandPaths()

	paths.MustMakeDir(paths.MustGetArtifactPath("archives"), permbits.MustString("ug=rwx,o=rx"))

	for _, cmdPath := range pathList {
		ct := helper.NewCommandTemplate(false, cmdPath)

		if err := buildPlatformIterator(ctx, ct, func(_ context.Context, ct *helper.CommandTemplate) error {
			// buildArtifact(ctx, ct)
			loga.PrintInfo("Create Archive for %s/%s (%s)",
				ct.GoOS, ct.GoArch,
				ct.OutputArtifact,
			)

			if err := buildTar(ct); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func buildTar(ct *helper.CommandTemplate) error {
	fileName := fmt.Sprintf("%s_%s_%s.tar.gz", ct.CommandName, ct.GoOS, ct.GoArch)
	if ct.GoArm != "" {
		fileName = fmt.Sprintf("%s_%s_%s_%s.tar.gz", ct.CommandName, ct.GoOS, ct.GoArch, ct.GoArm)
	}

	archivePath := paths.MustGetArtifactPath("archives", fileName)

	var f *os.File
	{
		var err error
		f, err = os.Create(archivePath)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	gzipWriter := gzip.NewWriter(f)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	if readmeFile := paths.MustGetWD("README.md"); paths.FileExists(readmeFile) {
		if err := tarAddFile(tarWriter, readmeFile, permbits.MustString("ugo=r")); err != nil {
			return err
		}
	}

	if binFile := ct.OutputArtifact; paths.FileExists(binFile) {
		if err := tarAddFile(tarWriter, binFile, permbits.MustString("ug=rwx,o=rx")); err != nil {
			return err
		}
	}

	return nil
}

func tarAddFile(tarWriter *tar.Writer, filename string, mode os.FileMode) error {
	var src *os.File
	{
		var err error
		src, err = os.Open(filename)
		if err != nil {
			return fmt.Errorf("unable to open source file[%s]: %w", filename, err)
		}
	}

	var srcStat os.FileInfo
	{
		var err error
		srcStat, err = src.Stat()
		if err != nil {
			return fmt.Errorf("unable to stat source file[%s]: %w", filename, err)
		}
	}

	hdr := &tar.Header{
		Name:    filepath.Base(filename),
		Mode:    int64(mode),
		Size:    srcStat.Size(),
		ModTime: build.GitCommitTime(),
	}

	if err := tarWriter.WriteHeader(hdr); err != nil {
		return fmt.Errorf("unable to write tar header for file[%s]: %w", filename, err)
	}

	if _, err := io.Copy(tarWriter, src); err != nil {
		return fmt.Errorf("unable to write contents to tar for file[%s]: %w", filename, err)
	}

	return nil
}
