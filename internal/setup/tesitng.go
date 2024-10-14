package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func SetupTesting(projectPath, projectType string) error {
	testDir := filepath.Join(projectPath, "tests")
	err := os.MkdirAll(testDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create tests directory: %v", err)
	}

	testFilePath := filepath.Join(testDir, "sample_test."+getFileExtension(projectType))
	testContent := getTestContent(projectType)
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create sample test file: %v", err)
	}

	color.Green("Testing framework set up successfully.")
	return nil
}

func getTestContent(projectType string) string {
	switch projectType {
	case "Go":
		return `package main

import "testing"

func TestSample(t *testing.T) {
	// Add your test here
}`
	case "Next.js":
		return `import { render, screen } from '@testing-library/react'
import Home from '../pages/index'

describe('Home', () => {
  it('renders a heading', () => {
    render(<Home />)
    const heading = screen.getByRole('heading', { level: 1 })
    expect(heading).toBeInTheDocument()
  })
})`
	// Add more project types as needed
	default:
		return "# Add your test content here"
	}
}

func getFileExtension(projectType string) string {
	switch projectType {
	case "Go":
		return "go"
	case "Next.js":
		return "js"
	// Add more project types as needed
	default:
		return "txt"
	}
}
