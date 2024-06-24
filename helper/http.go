package helper

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

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
	if eTag != nil && FileExists(filename) {
		eTagVal = eTag.Value
	}

	var resp *resty.Response
	{
		var err error
		PrintDebug("Downloading %s", rawURL)
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

	PrintDebug("[%s] Status Code: %d", rawURL, resp.StatusCode())

	switch resp.StatusCode() {
	case http.StatusNotModified:
		if err := os.Remove(filename + ".tmp"); err != nil {
			return fmt.Errorf("unable to remove temporary file: %w", err)
		}
	case http.StatusOK:
		if eTag != nil {
			eTag.Value = resp.Header().Get("etag")
			if err := eTag.Save(); err != nil {
				PrintDebug("Unable to write etag file: %s", err)
				return err
			}
		}
		if err := os.Rename(filename+".tmp", filename); err != nil {
			PrintDebug("Unable to rename from tmp file: %s", err)
			return fmt.Errorf("unable to replace file with temporary file: %w", err)
		}
	}

	if fileperm != 0 {
		if err := os.Chmod(filename, fileperm); err != nil {
			PrintDebug("Unable to change mode on file: %s", err)
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
	PrintDebug("Response Info:")
	PrintDebug("  Error      : %s", err)
	PrintDebug("  Status Code: %d", resp.StatusCode())
	PrintDebug("  Status     : %s", resp.Status())
	PrintDebug("  Proto      : %s", resp.Proto())
	PrintDebug("  Time       : %s", resp.Time())
	PrintDebug("  Received At: %s", resp.ReceivedAt())
	PrintDebug("  Headers    : %s", resp.Header())
	PrintDebug("  Body       : %s\n", resp)
	PrintDebug("")

	PrintDebug("Request Info:")
	PrintDebug("  URL        : %s", resp.Request.URL)
	PrintDebug("  Headers    : %s", resp.Request.Header)
	PrintDebug("")

	// Explore trace info
	PrintDebug("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	PrintDebug("  DNSLookup     : %s", ti.DNSLookup)
	PrintDebug("  ConnTime      : %s", ti.ConnTime)
	PrintDebug("  TCPConnTime   : %s", ti.TCPConnTime)
	PrintDebug("  TLSHandshake  : %s", ti.TLSHandshake)
	PrintDebug("  ServerTime    : %s", ti.ServerTime)
	PrintDebug("  ResponseTime  : %s", ti.ResponseTime)
	PrintDebug("  TotalTime     : %s", ti.TotalTime)
	PrintDebug("  IsConnReused  : %t", ti.IsConnReused)
	PrintDebug("  IsConnWasIdle : %t", ti.IsConnWasIdle)
	PrintDebug("  ConnIdleTime  : %s", ti.ConnIdleTime)
	PrintDebug("  RequestAttempt: %d", ti.RequestAttempt)
	PrintDebug("  RemoteAddr    : %s", ti.RemoteAddr.String())
	PrintDebug("")
}
