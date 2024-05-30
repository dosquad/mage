package helper

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
)

const (
	copyInChunksSize = 1024
)

func ExtractArchive(src, dest string) error {
	PrintDebug("Extract Archive: %s", src)
	if _, statErr := os.Stat(dest); os.IsNotExist(statErr) {
		PrintDebug("Create Directory: %s", dest)
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
	PrintDebug("URL: %s", u)

	if u.Scheme == "http" || u.Scheme == "https" {
		destArchive, err := DownloadToCache(src)
		if err != nil {
			return err
		}

		defer PrintDebug("ExtractArchive finished")

		return Unzip(destArchive, dest)
	}

	return fmt.Errorf("unknown scheme on source (%s)", src)
}

func getFilenameForURL(src string) (string, error) {
	client := resty.New()
	var resp *resty.Response
	{
		var err error
		resp, err = client.R().Head(src)
		if err != nil {
			return "", err
		}
	}

	resp.Header().Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(resp.Header().Get("Content-Disposition"))
	if err != nil {
		return "", err
	}
	filename := params["filename"]

	return filepath.Base(filename), nil
}

func DownloadToCache(src string) (string, error) {
	{
		filename, err := getFilenameForURL(src)
		if err != nil {
			return "", err
		}

		if v := filepath.Join(mg.CacheDir(), filename); FileExists(v) {
			return v, nil
		}
	}

	var tmpFile string
	{
		f, err := os.CreateTemp(mg.CacheDir(), "download-archive*")
		if err != nil {
			return "", err
		}
		tmpFile = f.Name()
		// if len(tmpDir) < len("download-archive") {
		// 	return "", fmt.Errorf("temporary directory name too short: %s", tmpDir)
		// }
	}
	defer func() { _ = os.Remove(tmpFile) }()
	PrintDebug("Temporary File: %s", tmpFile)

	var resp *resty.Response
	{
		var err error
		client := resty.New()
		resp, err = client.R().
			SetOutput(tmpFile).
			Get(src)
		if err != nil {
			return "", err
		}
	}

	resp.Header().Get("Content-Disposition")
	var destArchive string
	{
		_, params, err := mime.ParseMediaType(resp.Header().Get("Content-Disposition"))
		if err != nil {
			return "", err
		}
		destArchive, err = SanitizeArchivePath(mg.CacheDir(), params["filename"])
		if err != nil {
			return "", err
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

// Sanitize archive file pathing from "G305: Zip Slip vulnerability".
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

func Unzip(src, dest string) error {
	PrintDebug("Unzip(%s, %s)", src, dest)
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

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		if err := extractAndWriteFile(dest, f); err != nil {
			return err
		}
	}

	return nil
}
