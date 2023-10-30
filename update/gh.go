package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"github.com/cheggaaa/pb/v3"
	"github.com/google/go-github/v30/github"
	"github.com/mstxq17/MoreFind/errx"
	"golang.org/x/net/context"
	"io"
	"io/fs"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	Owner    = "mstxq17"
	ToolName = "MoreFind"
)

var (
	extIfFound = ".exe"
	// 下载时间应该设置大一些防止网速不好的情况
	GlobalTimeout     = time.Duration(60) * time.Second
	DefaultHttpClient *http.Client
)

type AssetFileCallback func(path string, fileInfo fs.FileInfo, data io.Reader) error

type GHReleaseDownloader struct {
	owner      string
	repoName   string // by default repoName is toolName
	assetName  string
	AssetID    int
	Format     AssetFormat
	Latest     *github.RepositoryRelease
	ghClient   *github.Client
	httpClient *http.Client
}

// 获取最新发新版信息并报告错误
func (d *GHReleaseDownloader) getLatestRelease() error {
	release, resp, err := d.ghClient.Repositories.GetLatestRelease(context.Background(), d.owner, d.repoName)
	var rateLimitErr *github.RateLimitError
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return errx.NewMsgf("updater -> repo %v/%v not found got %v", d.owner, d.repoName)
		} else if errx.As(err, &rateLimitErr) {
			return errx.NewMsg("hit github ratelimit while downloading latest release")
		}
		if resp == nil {
			return errx.NewWrapError(err, "updater -> network connect error")
		}
		return errx.NewWrapError(err, "updater -> unknown error that not be handled")
	}
	d.Latest = release
	return nil
}

func NewghReleaseDownloader(RepoName string) (*GHReleaseDownloader, error) {
	var owner, repoName string
	if strings.Contains(RepoName, "/") {
		// if it has diagonal that means mstxq17/MoreFind
		// 如果/存在，则说明是 mstxq17/MoreFind 的形式
		arr := strings.Split(RepoName, "/")
		if len(arr) != 2 {
			return nil, errx.NewMsgf("update RepoName: %v cannot be parsed", RepoName)
		}
		owner = arr[0]
		repoName = arr[1]
	} else {
		owner = Owner
		repoName = RepoName
	}
	if repoName == "" {
		repoName = ToolName
		//return nil, errx.NewMsg("update RepoName: repoName name cannot be empty")
	}
	// 全局超时配置client
	httpClient := DefaultHttpClient
	ghrd := GHReleaseDownloader{ghClient: github.NewClient(httpClient), repoName: repoName, owner: owner, httpClient: httpClient}
	err := ghrd.getLatestRelease()
	return &ghrd, err
}

func (d *GHReleaseDownloader) getToolAssetID(latest *github.RepositoryRelease) error {
	// MoreFind_1.4.6_darwin_arm64.tar.gz
	builder := &strings.Builder{}
	builder.WriteString(d.repoName)
	builder.WriteString("_")
	builder.WriteString(latest.GetTagName())
	builder.WriteString("_")
	builder.WriteString(runtime.GOOS)
	builder.WriteString("_")
	if strings.EqualFold(runtime.GOARCH, "amd64") {
		builder.WriteString("x86_64")
	} else {
		builder.WriteString(runtime.GOARCH)
	}
loop:
	for _, v := range latest.Assets {
		asset := v.GetName()
		switch {
		case strings.Contains(asset, Tar.FileExtension()):
			if strings.EqualFold(asset, builder.String()+Tar.FileExtension()) {
				d.AssetID = int(v.GetID())
				d.Format = Tar
				d.assetName = asset
				break loop
			}
		case strings.Contains(asset, Zip.FileExtension()):
			if strings.EqualFold(asset, builder.String()+Zip.FileExtension()) {
				d.AssetID = int(v.GetID())
				d.Format = Zip
				d.assetName = asset
				break loop
			}
		}
	}
	builder.Reset()
	// handle if id is zero (no asset found)
	if d.AssetID == 0 {
		return errx.NewMsgf("updater: could not find release asset for your platform (%s/%s)", runtime.GOOS, runtime.GOARCH)
	}
	return nil
}

// downloadAssetwithID
func (d *GHReleaseDownloader) downloadAssetwithID(id int64) (*http.Response, error) {
	_, rdurl, err := d.ghClient.Repositories.DownloadReleaseAsset(context.Background(), d.owner, d.repoName, id, nil)
	if err != nil {
		return nil, err
	}
	resp, err := d.httpClient.Get(rdurl)
	if err != nil {
		return nil, errx.NewMsg("failed to download release asset")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errx.NewMsgf("something went wrong got %v while downloading asset, expected status 200", resp.StatusCode)
	}
	if resp.Body == nil {
		return nil, errx.NewMsg("something went wrong got response without body")
	}
	return resp, nil
}

func (d *GHReleaseDownloader) DownloadTool() (*bytes.Buffer, error) {
	if err := d.getToolAssetID(d.Latest); err != nil {
		return nil, err
	}
	resp, err := d.downloadAssetwithID(int64(d.AssetID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !HideProgressBar {
		bar := pb.New64(resp.ContentLength).SetMaxWidth(100)
		bar.Start()
		resp.Body = bar.NewProxyReader(resp.Body)
		defer bar.Finish()
	}
	bin, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errx.NewMsg("failed to read response body")
	}
	return bytes.NewBuffer(bin), nil
}

func (d *GHReleaseDownloader) GetExecutableFromAsset() ([]byte, error) {
	var bin []byte
	var err error
	getToolCallback := func(path string, fileInfo fs.FileInfo, data io.Reader) error {
		if !strings.EqualFold(strings.TrimSuffix(fileInfo.Name(), extIfFound), ToolName) {
			return nil
		}
		bin, err = io.ReadAll(data)
		return err
	}
	buff, err := d.DownloadTool()
	if err != nil {
		return nil, err
	}
	_ = UnpackAssetWithCallback(d.Format, bytes.NewReader(buff.Bytes()), getToolCallback)
	return bin, errx.NewWrapError(err, "executable not found in archive") // Note: WrapfWithNil wraps msg if err != nil
}

// UnpackAssetWithCallback unpacks asset and executes callback function on every file in data
func UnpackAssetWithCallback(format AssetFormat, data *bytes.Reader, callback AssetFileCallback) error {
	if format != Zip && format != Tar {
		return errx.NewMsg("unpack -> github asset format not supported. only zip and tar are supported")
	}
	if format == Zip {
		zipReader, err := zip.NewReader(data, data.Size())
		if err != nil {
			return err
		}
		for _, f := range zipReader.File {
			data, err := f.Open()
			if err != nil {
				return err
			}
			if err := callback(f.Name, f.FileInfo(), data); err != nil {
				return err
			}
			_ = data.Close()
		}
	} else if format == Tar {
		gzipReader, err := gzip.NewReader(data)
		if err != nil {
			return err
		}
		tarReader := tar.NewReader(gzipReader)
		// iterate through the files in the archive
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if err := callback(header.Name, header.FileInfo(), tarReader); err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	DefaultHttpClient = &http.Client{
		Timeout: GlobalTimeout,
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
