package helper

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

func HTTPWriteFile(rawURL, filename string, eTag *ETagItem, fileperm os.FileMode) error {
	client := resty.New()
	eTagVal := ""
	if eTag != nil && FileExists(filename) {
		eTagVal = eTag.Value
	}

	var resp *resty.Response
	{
		var err error
		resp, err = client.R().
			SetOutput(filename+".tmp").
			SetHeader("If-None-Match", eTagVal).
			EnableTrace().
			Get(rawURL)
		if err != nil {
			return err
		}
	}

	switch resp.StatusCode() {
	case http.StatusNotModified:
		if err := os.Remove(filename + ".tmp"); err != nil {
			return fmt.Errorf("unable to remove temporary file: %w", err)
		}
	case http.StatusOK:
		if eTag != nil {
			eTag.Value = resp.Header().Get("etag")
			if err := eTag.Save(); err != nil {
				return err
			}
		}
		if err := os.Rename(filename+".tmp", filename); err != nil {
			return fmt.Errorf("unable to replace file with temporary file: %w", err)
		}
	}

	if fileperm != 0 {
		if err := os.Chmod(filename, fileperm); err != nil {
			return err
		}
	}

	return nil
}

func HTTPGetLatestGitHubVersion(slug string) (string, error) {
	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	resp, _ := client.R().
		Head(fmt.Sprintf("https://github.com/%s/releases/latest", slug))

	if resp == nil {
		return "", errors.New("response is nil")
	}

	location := resp.Header().Get("location")
	if strings.Contains(location, "/releases/tag/") {
		sp := strings.Split(location, "/releases/tag/")
		// return strings.TrimPrefix(sp[len(sp)-1], "v"), nil
		return sp[len(sp)-1], nil
	}

	return "", errors.New("unable to parse location")
}
