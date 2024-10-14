package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
)

func SetupGit(projectPath, projectType string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %v", err)
	}

	gitignorePath := filepath.Join(projectPath, ".gitignore")
	gitignoreContent := GetGitignoreContent(projectType)
	err = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create .gitignore file: %v", err)
	}

	color.Green("Git repository initialized and .gitignore created successfully.")
	return nil
}

func GetGitignoreContent(projectType string) string {
	switch projectType {
	case "Go":
		return `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/`
	case "Next.js":
		return `# Dependencies
/node_modules
/.pnp
.pnp.js

# Testing
/coverage

# Next.js
/.next/
/out/

# Production
/build

# Misc
.DS_Store
*.pem

# Debug
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Local env files
.env.local
.env.development.local
.env.test.local
.env.production.local

# Vercel
.vercel`
	// Add more project types as needed
	default:
		return "# Add your .gitignore content here"
	}
}
