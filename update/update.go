package update

import (
	"bytes"
	"github.com/Masterminds/semver/v3"
	"github.com/minio/selfupdate"
	"github.com/mstxq17/MoreFind/errx"
	"log"
	"os"
)

var (
	HideProgressBar  = true
	HideReleaseNotes = true
)

type ErrorCallback func() *log.Logger

// GetUpdateToolCallback reserve over-design for study
// 保留一定的冗余的"SB"设计用来学习
func GetUpdateToolCallback(toolName, version string, errorCallback ErrorCallback) func() {
	return GetUpdateToolFromRepoCallback(toolName, version, "", errorCallback)
}

func GetUpdateToolFromRepoCallback(toolName, version, repoName string, errorCallback ErrorCallback) func() {
	return func() {
		logger := errorCallback()
		if repoName == "" {
			repoName = toolName
		}
		gh, err := NewghReleaseDownloader(repoName)
		if err != nil {
			logger.Fatal(err)
		}
		latestVersion, err := semver.NewVersion(gh.Latest.GetTagName())
		currentVersion, err := semver.NewVersion(version)
		if err != nil {
			logger.Fatal(errx.NewWithMsgf(err, "failed to parse semversion from tagname `%v` got %v", gh.Latest.GetTagName()))
		}
		logger.Printf("Get Latest Version: v%v", latestVersion.String())
		if !IsOutdated(currentVersion.String(), latestVersion.String()) {
			logger.Printf("%v is already updated to the latest version: v%v", toolName, latestVersion.String())
			os.Exit(0)
		}
		// check permissions before downloading release
		updateOpts := selfupdate.Options{}
		if err := updateOpts.CheckPermissions(); err != nil {
			logger.Fatal(errx.NewWithMsgf(err, "update of %v %v -> %v failed , insufficient permission detected got: %v", toolName, currentVersion.String(), latestVersion.String()))
		}
		HideProgressBar = false
		bin, err := gh.GetExecutableFromAsset()
		if err != nil {
			logger.Fatal(errx.NewWithMsgf(err, "executable %v not found in release assetID `%v` got", toolName, gh.AssetID))
		}
		if err = selfupdate.Apply(bytes.NewBuffer(bin), updateOpts); err != nil {
			logger.Printf("update of %v %v -> %v failed, rolling back update", toolName, currentVersion.String(), latestVersion.String())
			if err := selfupdate.RollbackError(err); err != nil {
				logger.Println("")
				logger.Printf("updater -> rollback of update of %v failed got %v,pls reinstall %v", toolName, err, toolName)
			}
			os.Exit(1)
		}
		logger.Printf("%v successfully updated %v -> %v (latest)", toolName, currentVersion.String(), latestVersion.String())
		if !HideReleaseNotes {
			output := gh.Latest.GetBody()
			logger.Printf("\n\n%v\n", output)
		}
		os.Exit(0)
	}
}

// IsOutdated returns true if current version is outdated
func IsOutdated(current, latest string) bool {
	currentVer, _ := semver.NewVersion(current)
	latestVer, _ := semver.NewVersion(latest)
	if currentVer == nil || latestVer == nil {
		return current != latest
	}
	return latestVer.GreaterThan(currentVer)
}
