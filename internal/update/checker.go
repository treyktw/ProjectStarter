package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

const versionURL = "https://api.github.com/repos/treyktw/ProjectStarter/releases/latest"

type VersionInfo struct {
	LatestVersion string `json:"tag_name"`
	DownloadURL   string `json:"html_url"`
	AssetsURL     string `json:"assets_url"`
}

func CheckForUpdates(currentVersion string) (*VersionInfo, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", versionURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "project-starter-cli")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var versionInfo VersionInfo
	if err := json.Unmarshal(body, &versionInfo); err != nil {
		return nil, fmt.Errorf("failed to parse version info: %v", err)
	}

	fmt.Printf("Latest version from GitHub: %s\n", versionInfo.LatestVersion)

	// Remove any prefix (like "cli-tool-v" or "v") from the version string
	versionInfo.LatestVersion = strings.TrimPrefix(versionInfo.LatestVersion, "cli-tool-v")
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

func CheckAndPromptForUpdate(currentVersion string) error {
	versionInfo, err := CheckForUpdates(currentVersion)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %v", err)
	}

	if versionInfo == nil {
		// No update available
		return nil
	}

	fmt.Printf("A new version is available: %s (you're on %s)\n", versionInfo.LatestVersion, currentVersion)
	fmt.Print("Do you want to update now? (y/n): ")

	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		return fmt.Errorf("failed to read user input: %v", err)
	}

	if strings.ToLower(response) == "y" {
		fmt.Println("Starting update process...")
		return SelfUpdate(currentVersion)
	}

	fmt.Println("Update skipped. You can update later by running 'project-starter update'")
	return nil
}
