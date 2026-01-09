package web

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/dosquad/mage/helper/paths"
	"github.com/dosquad/mage/loga"
	"github.com/go-resty/resty/v2"
	"github.com/magefile/mage/mg"
)

//nolint:gocognit // ignore complexity.
func HTTPWriteFile(rawURL, filename string, eTag *ETagItem, fileperm os.FileMode, opts ...RestyOpt) error {
	client := resty.New()
	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}

	eTagVal := ""
	if eTag != nil && paths.FileExists(filename) {
		eTagVal = eTag.Value
	}

	var resp *resty.Response
	{
		var err error
		loga.PrintDebugf("Downloading %s", rawURL)
		resp, err = client.R().
			SetOutput(filename+".tmp").
			SetHeader("If-None-Match", eTagVal).
			EnableTrace().
			Get(rawURL)
		if err != nil {
			return err
		}

		if mg.Verbose() || mg.Debug() {
			RestyTrace(resp, err)
		}
	}

	loga.PrintDebugf("[%s] Status Code: %d", rawURL, resp.StatusCode())

	switch resp.StatusCode() {
	case http.StatusNotModified:
		if err := os.Remove(filename + ".tmp"); err != nil {
			return fmt.Errorf("unable to remove temporary file: %w", err)
		}
	case http.StatusOK:
		if eTag != nil {
			eTag.Value = resp.Header().Get("etag")
			if err := eTag.Save(); err != nil {
				loga.PrintDebugf("Unable to write etag file: %s", err)
				return err
			}
		}
		if err := os.Rename(filename+".tmp", filename); err != nil {
			loga.PrintDebugf("Unable to rename from tmp file: %s", err)
			return fmt.Errorf("unable to replace file with temporary file: %w", err)
		}
	}

	if fileperm != 0 {
		if err := os.Chmod(filename, fileperm); err != nil {
			loga.PrintDebugf("Unable to change mode on file: %s", err)
			return err
		}
	}

	return nil
}

func HTTPGetLatestGitHubVersion(slug string, opts ...RestyOpt) (string, error) {
	client := resty.New()
	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}

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

// repositoryRelease represents a GitHub release in a repository.
type repositoryRelease struct {
	TagName string `json:"tag_name,omitempty"`
	Name    string `json:"name,omitempty"`
}

type repositoryReleaseResult []repositoryRelease

func HTTPGetLatestGitHubReleaseMatchingTag(slug string, r *regexp.Regexp, opts ...RestyOpt) (string, error) {
	client := resty.New()
	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	resp, _ := client.R().
		SetResult(repositoryReleaseResult{}).
		SetHeader("Content-Type", "application/vnd.github+json").
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		Get(fmt.Sprintf("https://api.github.com/repos/%s/releases", slug))

	if resp == nil {
		return "", errors.New("response is nil")
	}

	if v, ok := resp.Result().(*repositoryReleaseResult); ok && v != nil {
		for _, release := range *v {
			if r.MatchString(release.TagName) {
				return release.TagName, nil
			}
		}
	}

	return "", errors.New("matching tag not found")
}

func RestyTrace(resp *resty.Response, err error) {
	// Explore response object
	loga.PrintDebugf("Response Info:")
	loga.PrintDebugf("  Error      : %s", err)
	loga.PrintDebugf("  Status Code: %d", resp.StatusCode())
	loga.PrintDebugf("  Status     : %s", resp.Status())
	loga.PrintDebugf("  Proto      : %s", resp.Proto())
	loga.PrintDebugf("  Time       : %s", resp.Time())
	loga.PrintDebugf("  Received At: %s", resp.ReceivedAt())
	loga.PrintDebugf("  Headers    : %s", resp.Header())
	loga.PrintDebugf("  Body       : %s\n", resp)
	loga.PrintDebugf("")

	loga.PrintDebugf("Request Info:")
	loga.PrintDebugf("  URL        : %s", resp.Request.URL)
	loga.PrintDebugf("  Headers    : %s", resp.Request.Header)
	loga.PrintDebugf("")

	// Explore trace info
	loga.PrintDebugf("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	loga.PrintDebugf("  DNSLookup     : %s", ti.DNSLookup)
	loga.PrintDebugf("  ConnTime      : %s", ti.ConnTime)
	loga.PrintDebugf("  TCPConnTime   : %s", ti.TCPConnTime)
	loga.PrintDebugf("  TLSHandshake  : %s", ti.TLSHandshake)
	loga.PrintDebugf("  ServerTime    : %s", ti.ServerTime)
	loga.PrintDebugf("  ResponseTime  : %s", ti.ResponseTime)
	loga.PrintDebugf("  TotalTime     : %s", ti.TotalTime)
	loga.PrintDebugf("  IsConnReused  : %t", ti.IsConnReused)
	loga.PrintDebugf("  IsConnWasIdle : %t", ti.IsConnWasIdle)
	loga.PrintDebugf("  ConnIdleTime  : %s", ti.ConnIdleTime)
	loga.PrintDebugf("  RequestAttempt: %d", ti.RequestAttempt)
	loga.PrintDebugf("  RemoteAddr    : %s", ti.RemoteAddr.String())
	loga.PrintDebugf("")
}
