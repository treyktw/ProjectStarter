package project

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

func BackupProject(dirPath string) error {
	projects, err := GetDirectories(dirPath)
	if err != nil {
		return fmt.Errorf("error getting projects: %v", err)
	}

	var selectedProject string
	prompt := &survey.Select{
		Message:  "Select a project to backup:",
		Options:  projects,
		PageSize: 15,
	}

	err = survey.AskOne(prompt, &selectedProject)
	if err != nil {
		return fmt.Errorf("project selection failed: %v", err)
	}

	projectPath := filepath.Join(dirPath, selectedProject)
	backupPath := filepath.Join(dirPath, selectedProject+"_backup_"+time.Now().Format("20060102_150405")+".zip")

	// Create and start a new spinner
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " Preparing to backup project..."
	s.Start()

	// Count files for progress bar
	fileCount, err := CountFiles(projectPath)
	s.Stop()
	if err != nil {
		return fmt.Errorf("error counting files: %v", err)
	}

	// Create progress bar
	bar := progressbar.NewOptions(fileCount,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan]Backing up project...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	// Create zip file
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the project directory
	err = filepath.Walk(projectPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set compression
		header.Method = zip.Deflate

		// Set relative path
		relPath, err := filepath.Rel(projectPath, filePath)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		}

		// Create file in zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		bar.Add(1)
		return nil
	})

	if err != nil {
		return fmt.Errorf("error creating backup: %v", err)
	}

	color.Green("\nBackup created successfully: %s", backupPath)
	return nil
}

func CountFiles(dir string) (int, error) {
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

func GetDirectories(path string) ([]string, error) {
	var dirs []string
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs, nil
}
