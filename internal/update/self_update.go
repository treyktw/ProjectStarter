package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

func SelfUpdate(currentVersion string) error {
	versionInfo, err := CheckForUpdates(currentVersion)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %v", err)
	}

	if versionInfo == nil {
		fmt.Println("You're already on the latest version.")
		return nil
	}

	fmt.Printf("New version available: %s\n", versionInfo.LatestVersion)
	fmt.Println("Starting update process...")

	// Determine the correct asset to download based on the OS and architecture
	fmt.Printf("Fetching assets from: %s\n", versionInfo.AssetsURL)
	assetURL, err := getAssetURL(versionInfo.AssetsURL)
	if err != nil {
		return fmt.Errorf("failed to get asset URL: %v", err)
	}
	fmt.Printf("Found asset URL: %s\n", assetURL)

	// Download the new version
	resp, err := http.Get(assetURL)
	if err != nil {
		return fmt.Errorf("failed to download new version: %v", err)
	}
	defer resp.Body.Close()

	// Create a temporary file to store the downloaded binary
	tmpFile, err := os.CreateTemp("", "project-starter-update")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the downloaded content to the temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Get the current executable path
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %v", err)
	}

	// Rename the current executable (essentially backing it up)
	backupPath := exe + ".old"
	err = os.Rename(exe, backupPath)
	if err != nil {
		return fmt.Errorf("failed to rename current executable: %v", err)
	}

	// Move the new version to replace the current executable
	err = os.Rename(tmpFile.Name(), exe)
	if err != nil {
		// If moving the new version fails, restore the old version
		os.Rename(backupPath, exe)
		return fmt.Errorf("failed to move new version: %v", err)
	}

	// Remove the backup of the old version
	os.Remove(backupPath)

	fmt.Println("Update successful! Please restart the application.")
	return nil
}

func getAssetURL(assetsURL string) (string, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", assetsURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "project-starter-cli")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get assets: %v", err)
	}
	defer resp.Body.Close()

	var assets []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		return "", fmt.Errorf("failed to parse assets: %v", err)
	}

	osArch := runtime.GOOS + "-" + runtime.GOARCH
	for _, asset := range assets {
		assetName := strings.ToLower(asset.Name)
		if strings.Contains(assetName, runtime.GOOS) && strings.Contains(assetName, runtime.GOARCH) {
			return asset.BrowserDownloadURL, nil
		}
	}

	// If no exact match found, try to find a close match
	for _, asset := range assets {
		assetName := strings.ToLower(asset.Name)
		if strings.Contains(assetName, runtime.GOOS) {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no suitable asset found for %s", osArch)
}
