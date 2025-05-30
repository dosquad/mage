package web

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/loga"
	"github.com/go-resty/resty/v2"
	"github.com/h2non/filetype"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
)

const (
	copyInChunksSize = 1024
)

func ExtractArchive(src, dest string, opts ...RestyOpt) error {
	loga.PrintDebug("Extract Archive: %s", src)
	if _, statErr := os.Stat(dest); os.IsNotExist(statErr) {
		loga.PrintDebug("Create Directory: %s", dest)
		if err := os.MkdirAll(dest, permbits.MustString("ug=rwx,o=rx")); err != nil {
			return err
		}
	}

	var u *url.URL
	{
		var err error
		u, err = url.Parse(src)
		if err != nil {
			return err
		}
	}
	loga.PrintDebug("URL: %s", u)

	if u.Scheme == "http" || u.Scheme == "https" { // URL
		destArchive, err := DownloadToCache(src, opts...)
		if err != nil {
			return err
		}

		defer loga.PrintDebug("ExtractArchive finished")

		return Uncompress(destArchive, dest)
	}

	if strings.HasPrefix(src, "/") { // absolute path
		defer loga.PrintDebug("ExtractArchive finished")

		return Uncompress(src, dest)
	}

	return fmt.Errorf("unknown scheme on source (%s)", src)
}

func GetFilenameForURL(src string, opts ...RestyOpt) (string, error) {
	client := resty.New()
	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}

	var resp *resty.Response
	{
		var err error
		resp, err = client.R().Head(src)
		// if mg.Debug() {
		// 	RestyTrace(resp, err)
		// }
		if err != nil {
			return "", err
		}
	}

	{
		cdHeader := resp.Header().Get("Content-Disposition")
		if cdHeader != "" {
			_, params, err := mime.ParseMediaType(resp.Header().Get("Content-Disposition"))
			if err != nil {
				return "", err
			}
			return filepath.Base(params["filename"]), nil
		}
	}

	return filepath.Base(src), nil
}

func DownloadToPath(src, dest string, opts ...RestyOpt) (string, error) {
	loga.PrintDebug("DownloadToPath(src:%s, dest:%s opts...)", src, dest)
	{
		filename, err := GetFilenameForURL(src, opts...)
		if err != nil {
			return "", err
		}

		if v := filepath.Join(dest, filename); paths.FileExists(v) {
			return v, nil
		}
	}

	var tmpFile string
	{
		f, err := os.CreateTemp(dest, "download-archive*")
		if err != nil {
			return "", err
		}
		tmpFile = f.Name()
		// if len(tmpDir) < len("download-archive") {
		// 	return "", fmt.Errorf("temporary directory name too short: %s", tmpDir)
		// }
	}
	defer func() { _ = os.Remove(tmpFile) }()
	loga.PrintDebug("Temporary File: %s", tmpFile)

	var resp *resty.Response
	{
		var err error
		client := resty.New()
		for _, opt := range opts {
			if opt != nil {
				opt(client)
			}
		}
		resp, err = client.R().
			SetOutput(tmpFile).
			Get(src)
		if err != nil {
			return "", err
		}
	}

	var destArchive string
	{
		cdHeader := resp.Header().Get("Content-Disposition")
		if cdHeader != "" {
			_, params, err := mime.ParseMediaType(cdHeader)
			if err != nil {
				return "", err
			}
			destArchive, err = SanitizeArchivePath(dest, params["filename"])
			if err != nil {
				return "", err
			}
		} else {
			var err error
			destArchive, err = SanitizeArchivePath(dest, filepath.Base(src))
			if err != nil {
				return "", err
			}
		}
	}

	if err := os.Rename(
		tmpFile,
		destArchive,
	); err != nil {
		return "", err
	}

	return destArchive, nil
}

func DownloadToCache(src string, opts ...RestyOpt) (string, error) {
	loga.PrintDebug("DownloadToCache(src:%s, opts...)", src)
	return DownloadToPath(src, mg.CacheDir(), opts...)
}

// SanitizeArchivePath archive file pathing from "G305: Zip Slip vulnerability".
func SanitizeArchivePath(dest, target string) (string, error) {
	v := filepath.Join(dest, target)
	if strings.HasPrefix(v, filepath.Clean(dest)) {
		return v, nil
	}

	return "", fmt.Errorf("%s: %s", "content filepath is tainted", target)
}

// Closure to address file descriptors issue with all the deferred .Close() methods.
func extractAndWriteFile(dest string, zf *zip.File) error {
	var path string
	{
		var err error
		// Check for ZipSlip: https://snyk.io/research/zip-slip-vulnerability
		path, err = SanitizeArchivePath(dest, zf.Name)
		if err != nil {
			return err
		}
	}

	var rc io.ReadCloser
	{
		var err error
		rc, err = zf.Open()
		if err != nil {
			return err
		}
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	if zf.FileInfo().IsDir() {
		// create the directory and return.
		return os.MkdirAll(
			path,
			permbits.Force(zf.Mode(),
				permbits.UserExecute, permbits.GroupExecute,
			),
		)
	}

	mode := zf.Mode()
	mode = permbits.Force(mode, permbits.UserWrite, permbits.GroupWrite)
	if err := os.MkdirAll(
		filepath.Dir(path),
		permbits.Force(mode,
			permbits.UserExecute, permbits.GroupExecute,
		),
	); err != nil {
		return err
	}

	var f *os.File
	{
		var err error
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
		if err != nil {
			return err
		}
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if err := copyInChunks(f, rc); err != nil {
		return err
	}

	return nil
}

func copyInChunks(dst io.Writer, src io.Reader) error {
	// totalRead := int64(0)
	for {
		_, err := io.CopyN(dst, src, copyInChunksSize)
		// totalRead += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}
	}

	return nil
}

func Uncompress(src, dest string) error {
	var buf []byte
	{
		f, fErr := os.Open(src)
		if fErr != nil {
			return fmt.Errorf("unable to open source[%s]: %w", src, fErr)
		}

		buf = make([]byte, 256) //nolint:mnd // only need first 256 bytes for identification.

		if _, err := f.Read(buf); err != nil {
			return fmt.Errorf("unable to read first 256 bytes[%s]: %w", src, err)
		}
	}
	// buf, _ := os.ReadFile(src)
	kind, err := filetype.Match(buf)
	if err != nil {
		return fmt.Errorf("unable to determine compression type: %w", err)
	}

	if kind == filetype.Unknown {
		return errors.New("unknown file type")
	}

	switch kind.MIME.Type { //nolint:gocritic // for future code expansion.
	case "application":
		switch kind.MIME.Subtype {
		case "gzip":
			return Untargz(src, dest)
		case "zip":
			return Unzip(src, dest)
		}
	}

	return nil
}

//nolint:gocognit // untar+decompress(gzip).
func Untargz(src, dest string) error {
	loga.PrintDebug("Untargz(%s, %s)", src, dest)
	dest = filepath.Clean(dest) + string(os.PathSeparator)

	var srcStream io.ReadCloser
	{
		var err error
		srcStream, err = os.Open(src)
		if err != nil {
			return fmt.Errorf("unable to open source file: %w", err)
		}
	}

	var tarStream *gzip.Reader
	{
		var err error
		tarStream, err = gzip.NewReader(srcStream)
		if err != nil {
			return fmt.Errorf("unable to decompress stream: %w", err)
		}
	}

	tarReader := tar.NewReader(tarStream)

	for {
		var header *tar.Header
		{
			var err error
			header, err = tarReader.Next()

			if err == io.EOF {
				break
			}

			if err != nil {
				return fmt.Errorf("untar failed: %w", err)
			}
		}

		switch header.Typeflag {
		case tar.TypeDir:
			var target string
			{
				var err error
				target, err = SanitizeArchivePath(dest, header.Name)
				if err != nil {
					return fmt.Errorf("sanitize path failed[%s]: %w", header.Name, err)
				}
			}

			if err := os.Mkdir(target, fs.FileMode(header.Mode)); err != nil { //nolint:gosec // G115 overflow possible.
				return fmt.Errorf("mkdir failed[%s]: %w", target, err)
			}
		case tar.TypeReg:
			var target string
			{
				var err error
				target, err = SanitizeArchivePath(dest, header.Name)
				if err != nil {
					return fmt.Errorf("sanitize path failed[%s]: %w", header.Name, err)
				}
			}

			var outFile *os.File
			{
				var err error
				outFile, err = os.Create(target)
				if err != nil {
					return fmt.Errorf("unable to create file[%s]: %w", target, err)
				}
			}

			if _, err := copyBlocks(outFile, tarReader); err != nil {
				return fmt.Errorf("unable to write file[%s]: %w", target, err)
			}

			if err := outFile.Close(); err != nil {
				return fmt.Errorf("unable to close file[%s]: %w", target, err)
			}
		default:
			return fmt.Errorf("unknown record type in tarball[%s]: type is %d", header.Name, header.Typeflag)
		}
	}

	return nil
}

func copyBlocks(dst io.Writer, src io.Reader) (int64, error) {
	var total int64

	for {
		n, err := io.CopyN(dst, src, 1024) //nolint:mnd // 1Kb blocks.
		total += n
		if err != nil {
			if errors.Is(err, io.EOF) {
				return total, nil
			}
			return total, err
		}
	}
}

func Unzip(src, dest string) error {
	loga.PrintDebug("Unzip(%s, %s)", src, dest)
	dest = filepath.Clean(dest) + string(os.PathSeparator)

	var r *zip.ReadCloser
	{
		var err error
		r, err = zip.OpenReader(src)
		if err != nil {
			return err
		}
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	if err := os.MkdirAll(dest, permbits.MustString("u=rwx,go=rw")); err != nil {
		return err
	}

	for _, f := range r.File {
		if err := extractAndWriteFile(dest, f); err != nil {
			return err
		}
	}

	return nil
}
