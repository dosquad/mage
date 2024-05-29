package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

func HTTPWriteFile(rawURL, filename string, fileperm os.FileMode) error {
	client := resty.New()
	if _, err := client.R().
		SetOutput(filename).
		Get(rawURL); err != nil {
		return err
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
		return strings.TrimPrefix(sp[len(sp)-1], "v"), nil
	}

	return "", errors.New("unable to parse location")
}
