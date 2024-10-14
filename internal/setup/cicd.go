package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func SetupCICD(projectPath, projectType string) error {
	cicdDir := filepath.Join(projectPath, ".github", "workflows")
	err := os.MkdirAll(cicdDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create .github/workflows directory: %v", err)
	}

	cicdFilePath := filepath.Join(cicdDir, "ci-cd.yml")
	cicdContent := getCICDContent(projectType)
	err = os.WriteFile(cicdFilePath, []byte(cicdContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create CI/CD configuration file: %v", err)
	}

	color.Green("CI/CD template added successfully.")
	return nil
}

func getCICDContent(projectType string) string {
	switch projectType {
	case "Go":
		return `name: Go CI/CD

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.16
    - name: Build
      run: go build -v ./...
    - name: Test
      run: go test -v ./...`
	case "Next.js":
		return `name: Next.js CI/CD

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Use Node.js
      uses: actions/setup-node@v2
      with:
        node-version: '14.x'
    - run: npm ci
    - run: npm run build
    - run: npm test`
	// Add more project types as needed
	default:
		return "# Add your CI/CD configuration here"
	}
}
