package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"

	"project-starter/internal/project"
	"project-starter/internal/update"
)

const Version = "1.0.0"

func main() {
	// Set up context for graceful exit
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	versionInfo, err := update.CheckForUpdates(Version)
	if err != nil {
		color.Yellow("Failed to check for updates: %v", err)
	} else if versionInfo != nil {
		color.Yellow("A new version is available: %s", versionInfo.LatestVersion)
		color.Yellow("Download it from: %s", versionInfo.DownloadURL)
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt in a separate goroutine
	go func() {
		<-interrupt
		fmt.Println("\n")
		color.Yellow("Operation canceled. Goodbye!")
		cancel()
		os.Exit(0)
	}()

	// Display welcome message
	welcomeFigure := figure.NewFigure("Welcome back", "", true)
	color.Red(welcomeFigure.String())
	color.Cyan("Trey")

	// Ask what to build
	color.Yellow("\nWhat would you like to build today?")

	// Start from D:/vscode/
	currentPath := "D:/vscode/"

	// Run the main loop with context
	err = run(ctx, currentPath)

	// Check if the context was canceled (i.e., Ctrl-C was pressed)
	if err != nil {
		if err == context.Canceled {
			fmt.Println("\n")
			color.Yellow("Operation canceled. Goodbye!")
		} else {
			color.Red("An error occurred: %v", err)
		}
	}
}

func run(ctx context.Context, startPath string) error {
	currentPath := startPath

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			dirs, err := project.GetDirectories(currentPath)
			if err != nil {
				return fmt.Errorf("error getting directories: %v", err)
			}

			dirs = append([]string{
				"[Use this directory]",
				"[Go back]",
				"[View Project Statistics]",
				"[Backup Project]",
			}, dirs...)

			var selected string
			prompt := &survey.Select{
				Message:  fmt.Sprintf("Current directory: %s\nSelect directory or action:", currentPath),
				Options:  dirs,
				PageSize: 15,
			}

			err = survey.AskOne(prompt, &selected)

			if err != nil {
				return fmt.Errorf("prompt failed: %v", err)
			}

			// Check context after user input
			if ctx.Err() != nil {
				return ctx.Err()
			}

			switch selected {
			case "[Use this directory]":
				return project.CreateProject(ctx, currentPath)
			case "[Go back]":
				currentPath = filepath.Dir(currentPath)
			case "[View Project Statistics]":
				if err := project.ViewProjectStatistics(currentPath); err != nil {
					if err == context.Canceled {
						return err
					}
					color.Red("Error viewing project statistics: %v", err)
				}
			case "[Backup Project]":
				if err := project.BackupProject(currentPath); err != nil {
					if err == context.Canceled {
						return err
					}
					color.Red("Error backing up project: %v", err)
				}
			default:
				currentPath = filepath.Join(currentPath, selected)
			}
		}
	}
}
