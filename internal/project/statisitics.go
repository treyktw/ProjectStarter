package project

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

type ProjectStats struct {
	LastModified time.Time
	TotalSize    int64
	FileCount    int
}

func ViewProjectStatistics(dirPath string) error {
	projects, err := GetDirectories(dirPath)
	if err != nil {
		return fmt.Errorf("error getting projects: %v", err)
	}

	var selectedProject string
	prompt := &survey.Select{
		Message:  "Select a project to view statistics:",
		Options:  projects,
		PageSize: 15,
	}

	err = survey.AskOne(prompt, &selectedProject)
	if err != nil {
		return fmt.Errorf("project selection failed: %v", err)
	}

	projectPath := filepath.Join(dirPath, selectedProject)

	// Create and start a new spinner
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " Gathering project statistics..."
	s.Start()

	stats, err := getProjectStats(projectPath)
	s.Stop()

	if err != nil {
		return fmt.Errorf("error getting project statistics: %v", err)
	}

	displayProjectStats(selectedProject, stats)
	return nil
}

func displayProjectStats(projectName string, stats ProjectStats) {
	color.Cyan("\nProject Statistics for %s:", projectName)
	color.Yellow("Last Modified: %s", stats.LastModified.Format("2006-01-02 15:04:05"))
	color.Yellow("Total Size: %s", humanize.Bytes(uint64(stats.TotalSize)))
	color.Yellow("Number of Files: %d", stats.FileCount)
}

func getProjectStats(projectPath string) (ProjectStats, error) {
	var stats ProjectStats
	var err error

	stats.LastModified, err = getLastModifiedTime(projectPath)
	if err != nil {
		return stats, err
	}

	// First, count the total number of files
	totalFiles := 0
	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalFiles++
		}
		return nil
	})
	if err != nil {
		return stats, err
	}

	// Create a progress bar
	bar := progressbar.NewOptions(totalFiles,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/2][reset] Counting files..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	// Now, gather the statistics with the progress bar
	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			stats.FileCount++
			stats.TotalSize += info.Size()
			bar.Add(1)
		}
		return nil
	})

	// Update progress bar description for size calculation
	bar.Describe("[cyan][2/2][reset] Calculating total size...")
	bar.Add(totalFiles - stats.FileCount) // Complete the bar

	return stats, err
}

func getLastModifiedTime(dirPath string) (time.Time, error) {
	var lastModTime time.Time
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.ModTime().After(lastModTime) {
			lastModTime = info.ModTime()
		}
		return nil
	})
	return lastModTime, err
}
