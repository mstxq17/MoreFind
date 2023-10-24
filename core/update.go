package core

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v30/github"
	"github.com/pkg/errors"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	Author          = "mstxq17"
	CheckPoint      = "https://raw.githubusercontent.com/%v/%v/master/version"
	ProxyCheckPoint = "https://ghproxy.com/%v"
)

var (
	HideProgressBar       = false
	HideReleaseNotes      = false
	VersionCheckTimeout   = time.Duration(5) * time.Second
	DownloadUpdateTimeout = time.Duration(5) * time.Second
	DefaultHttpClient     *http.Client
)

type GHReleaseDownloader struct {
	assetName     string // required assetName given as input
	fullAssetName string // full asset name of asset that contains tool for this platform
	organization  string // organization name of repo
	repoName      string // 项目仓库名 MoreFind
	AssetID       int
	Latest        *github.RepositoryRelease
	client        *github.Client
	httpClient    *http.Client
}

func (d *GHReleaseDownloader) getLatestRelease() error {
	release, resp, err := d.client.Repositories.GetLatestRelease(context.Background(), d.organization, d.repoName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return NewError(err, fmt.Sprintf("repo %v/%v not found got %v", d.organization, d.repoName, d.repoName))
		} else if _, ok := err.(*github.RateLimitError); ok {
			return NewError(err, "hit github rate-limit while downloading latest release")
		} else if resp != nil && (resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized) {
			return NewError(err, "gh auth failed try unsetting GITHUB_TOKEN env variable")
		}
	}
	d.Latest = release
	return nil

}

func (d *GHReleaseDownloader) GetExecutableFromAsset() ([]byte, error) {
	//var bin []byte
	_, err := d.DownloadTool()
	if err != nil {
		return nil, err
	}
	test := make([]byte, 0)
	return test, nil
}

// DownloadTool downloads tool and returns bin data
func (d *GHReleaseDownloader) DownloadTool() (*bytes.Buffer, error) {
	if err := d.getToolAssetID(d.Latest); err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)
	buffer.WriteString("123")
	return buffer, nil
}

func (d *GHReleaseDownloader) getToolAssetID(latest *github.RepositoryRelease) error {
	// MoreFind_Darwin_arm64.tar.gz
	builder := &strings.Builder{}
	builder.WriteString(d.assetName)
	builder.WriteString("_")
	fmt.Println(runtime.GOOS)
	if strings.EqualFold(runtime.GOOS, "darwin") {
		builder.WriteString("Darwin")
	} else {
		builder.WriteString(runtime.GOOS)
	}
	builder.WriteString("_")
	builder.WriteString(runtime.GOARCH)
	fmt.Println(builder.String())
	//loop:
	for _, v := range latest.Assets {
		asset := v.GetName()
		fmt.Println(asset)
	}
	return nil
}

func NewghReleaseDownloader(RepoName string) (*GHReleaseDownloader, error) {
	var orgName, repoName string
	if strings.Contains(RepoName, "/") {
		arr := strings.Split(RepoName, "/")
		if len(arr) != 2 {
			return nil, errors.New(fmt.Sprintf("update: invalid repo name %v", RepoName))
		}
		orgName = arr[0]
		repoName = arr[1]
	} else {
		orgName = Author
		repoName = RepoName
	}
	httpClient := &http.Client{
		Timeout: DownloadUpdateTimeout,
	}
	if orgName == "" {
		return nil, errors.New("update: organization name cannot be empty")
	}
	ghrd := GHReleaseDownloader{client: github.NewClient(httpClient), repoName: repoName, assetName: repoName, httpClient: httpClient, organization: orgName}
	err := ghrd.getLatestRelease()
	return &ghrd, err
}

func IsOutdated(current, latest string) bool {
	currentVer, _ := semver.NewVersion(current)
	latestVer, _ := semver.NewVersion(latest)
	if currentVer == nil || latestVer == nil {
		return currentVer != latestVer
	}
	return latestVer.GreaterThan(currentVer)
}

func init() {
	DefaultHttpClient = &http.Client{
		Timeout: VersionCheckTimeout,
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
