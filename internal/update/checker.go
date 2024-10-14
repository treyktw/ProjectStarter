package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	versionURL = "https://example.com/cli-version.json" // Replace with your actual URL
)

type VersionInfo struct {
	LatestVersion string `json:"latest_version"`
	DownloadURL   string `json:"download_url"`
}

func CheckForUpdates(currentVersion string) (*VersionInfo, error) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(versionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %v", err)
	}
	defer resp.Body.Close()

	var versionInfo VersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&versionInfo); err != nil {
		return nil, fmt.Errorf("failed to parse version info: %v", err)
	}

	current, err := semver.NewVersion(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid current version: %v", err)
	}

	latest, err := semver.NewVersion(versionInfo.LatestVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid latest version: %v", err)
	}

	if latest.GreaterThan(current) {
		return &versionInfo, nil
	}

	return nil, nil
}
