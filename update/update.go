package core

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cheggaaa/pb/v3"
	"github.com/google/go-github/v30/github"
	"github.com/pkg/errors"
	"io"
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
	Format        string
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
	test := make([]byte, 0)
	_, err := d.DownloadTool()
	if err != nil {
		return test, err
	}

	test = make([]byte, 0)
	return test, nil
}

// downloadAssetwithID
func (d *GHReleaseDownloader) downloadAssetwithID(id int64) (*http.Response, error) {
	_, rdurl, err := d.client.Repositories.DownloadReleaseAsset(context.Background(), d.organization, d.repoName, id, nil)
	if err != nil {
		return nil, err
	}
	resp, err := d.httpClient.Get(rdurl)
	if err != nil {
		return nil, errors.New("failed to download release asset")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("something went wrong got %v while downloading asset, expected status 200", resp.StatusCode))
	}
	if resp.Body == nil {
		return nil, errors.New("something went wrong got response without body")
	}
	return resp, nil
}

// DownloadTool downloads tool and returns bin data
func (d *GHReleaseDownloader) DownloadTool() (*bytes.Buffer, error) {
	if err := d.getToolAssetID(d.Latest); err != nil {
		return nil, err
	}
	resp, err := d.downloadAssetwithID(int64(d.AssetID))
	if err != nil {
		return nil, err
	}
	if !HideProgressBar {
		bar := pb.New64(resp.ContentLength).SetMaxWidth(100)
		bar.Start()
		resp.Body = bar.NewProxyReader(resp.Body)
		defer bar.Finish()
	}
	bin, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body")
	}
	fmt.Println(bytes.NewBuffer(bin).String())
	return bytes.NewBuffer(bin), nil
}

func (d *GHReleaseDownloader) getToolAssetID(latest *github.RepositoryRelease) error {
	// MoreFind_1.4.6_darwin_arm64.tar.gz
	builder := &strings.Builder{}
	builder.WriteString(d.assetName)
	builder.WriteString("_")
	builder.WriteString(latest.GetTagName())
	builder.WriteString("_")
	builder.WriteString(runtime.GOOS)
	builder.WriteString("_")
	builder.WriteString(runtime.GOARCH)
	tarAssetName := builder.String() + ".tar.gz"
	zipAssetName := builder.String() + ".zip"
loop:
	for _, v := range latest.Assets {
		asset := v.GetName()
		fmt.Println(zipAssetName, asset)
		switch {
		case strings.Contains(asset, "tar.gz"):
			if strings.EqualFold(asset, tarAssetName) {
				d.AssetID = int(v.GetID())
				d.Format = "tar.gz"
				d.fullAssetName = asset
				break loop
			}
		case strings.Contains(asset, "zip"):
			if strings.EqualFold(asset, zipAssetName) {
				d.AssetID = int(v.GetID())
				d.Format = "zip"
				d.fullAssetName = asset
				break loop
			}
		}
	}
	builder.Reset()
	// handle if id is zero (no asset found)
	if d.AssetID == 0 {
		return errors.New(fmt.Sprintf("update: could not find release asset for your platform (%s/%s)", runtime.GOOS, runtime.GOARCH))
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
