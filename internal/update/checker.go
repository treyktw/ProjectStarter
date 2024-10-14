package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	versionURL = "hhttps://github.com/treyktw/ProjectStarter/releases/tag/cli-tool" // Replace with your actual URL
)

type VersionInfo struct {
	LatestVersion string `json:"tag_name"`
	DownloadURL   string `json:"html_url"`
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
	// Remove the "v" prefix if your tags are like "v1.0.0"
	versionInfo.LatestVersion = strings.TrimPrefix(versionInfo.LatestVersion, "v")

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
