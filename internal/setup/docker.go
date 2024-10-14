package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func SetupDocker(projectPath, projectType string) error {
	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	dockerComposeFilePath := filepath.Join(projectPath, "docker-compose.yml")

	dockerfileContent := getDockerfileContent(projectType)
	err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create Dockerfile: %v", err)
	}

	dockerComposeContent := getDockerComposeContent(projectType)
	err = os.WriteFile(dockerComposeFilePath, []byte(dockerComposeContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create docker-compose.yml: %v", err)
	}

	color.Green("Docker support added successfully.")
	return nil
}

func getDockerfileContent(projectType string) string {
	switch projectType {
	case "Go":
		return `FROM golang:1.16-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]`
	case "Next.js":
		return `FROM node:14-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
CMD ["npm", "start"]`
	// Add more project types as needed
	default:
		return "# Add your Dockerfile content here"
	}
}

func getDockerComposeContent(projectType string) string {
	switch projectType {
	case "Go":
		return `version: '3'
services:
  app:
    build: .
    ports:
      - "8080:8080"`
	case "Next.js":
		return `version: '3'
services:
  app:
    build: .
    ports:
      - "3000:3000"`
	// Add more project types as needed
	default:
		return "# Add your docker-compose.yml content here"
	}
}
