package project

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"

	"project-starter/internal/setup"
)

type projectTemplate struct {
	name     string
	initFunc func(string, string) (*exec.Cmd, error)
}

var templates = []projectTemplate{
	{"Go", initGoProject},
	{"Next.js", initNextJSProject},
	{"Rust", initRustProject},
	{"Vite", initViteProject},
	{"Vue", initVueProject},
}

func CreateProject(_ context.Context, basePath string) error {
	var projectName string
	err := survey.AskOne(&survey.Input{Message: "Enter project name:"}, &projectName)
	if err != nil {
		return fmt.Errorf("project name input failed: %v", err)
	}

	projectPath := filepath.Join(basePath, projectName)

	err = os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating project directory: %v", err)
	}

	var result string
	prompt := &survey.Select{
		Message: "Select project type:",
		Options: getTemplateNames(),
	}
	err = survey.AskOne(prompt, &result)
	if err != nil {
		return fmt.Errorf("project type selection failed: %v", err)
	}

	var selectedTemplate projectTemplate
	for _, t := range templates {
		if t.name == result {
			selectedTemplate = t
			break
		}
	}

	cmd, err := selectedTemplate.initFunc(projectPath, projectName)
	if err != nil {
		color.Red("Error preparing project initialization: %v", err)
		if cmd == nil {
			color.Yellow("Please ensure the required tools are installed and available in your PATH.")
		}
		return err
	}

	if cmd != nil {
		cmd.Dir = projectPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		color.Cyan("Initializing project...")
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("error initializing project: %v", err)
		}
	}

	color.Green("Successfully created %s project in %s", result, projectPath)

	// Additional setup options
	setupOptions := []string{
		"Docker support",
		"CI/CD template",
		"Testing framework",
		"Git initialization",
	}

	var selectedOptions []string
	multiSelect := &survey.MultiSelect{
		Message: "Select additional setup options:",
		Options: setupOptions,
	}
	err = survey.AskOne(multiSelect, &selectedOptions)
	if err != nil {
		return fmt.Errorf("setup options selection failed: %v", err)
	}

	for _, option := range selectedOptions {
		switch option {
		case "Docker support":
			if err := setup.SetupDocker(projectPath, result); err != nil {
				color.Red("Error setting up Docker: %v", err)
			}
		case "CI/CD template":
			if err := setup.SetupCICD(projectPath, result); err != nil {
				color.Red("Error setting up CI/CD: %v", err)
			}
		case "Testing framework":
			if err := setup.SetupTesting(projectPath, result); err != nil {
				color.Red("Error setting up testing framework: %v", err)
			}
		case "Git initialization":
			if err := setup.SetupGit(projectPath, result); err != nil {
				color.Red("Error initializing Git: %v", err)
			}
		}
	}

	err = openInVSCode(projectPath)
	if err != nil {
		color.Red("Error opening project in VS Code: %v", err)
	} else {
		color.Green("Opened project in Visual Studio Code.")
	}

	return nil
}

func getTemplateNames() []string {
	names := make([]string, len(templates))
	for i, t := range templates {
		names[i] = t.name
	}
	return names
}

func openInVSCode(path string) error {
	cmd := exec.Command("code", ".")
	cmd.Dir = path
	return cmd.Run()
}

func initRustProject(path, projectName string) (*exec.Cmd, error) {
	cmd := exec.Command("cargo", "init")
	cmd.Dir = path
	return cmd, nil
}

func initViteProject(path, projectName string) (*exec.Cmd, error) {
	cmd := exec.Command("npm", "init", "vite@latest", projectName)
	cmd.Dir = filepath.Dir(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd, nil
}

func initVueProject(path, projectName string) (*exec.Cmd, error) {
	cmd := exec.Command("npm", "init", "vue@latest", projectName)
	cmd.Dir = filepath.Dir(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd, nil
}

func initNextJSProject(path, projectName string) (*exec.Cmd, error) {
	runtimes := []string{"npm", "pnpm", "bun", "deno"}

	runtimePrompt := promptui.Select{
		Label: "Select the runtime for your Next.js project",
		Items: runtimes,
	}

	_, runtime, err := runtimePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("runtime selection failed: %v", err)
	}

	if !isExecutableAvailable(runtime) {
		return nil, fmt.Errorf("the selected runtime '%s' is not available in your current PATH. Please install it or use a different terminal", runtime)
	}

	var cmd *exec.Cmd

	switch runtime {
	case "npm":
		cmd = exec.Command("npx", "create-next-app@latest", ".")
	case "pnpm":
		cmd = exec.Command("pnpm", "create", "next-app", ".")
	case "bun":
		cmd = exec.Command("bunx", "create-next-app", ".")
	case "deno":
		cmd = exec.Command("deno", "run", "--allow-env --allow-sys --allow-read --allow-write", "npm:create-next-app@latest", ".")
	default:
		return nil, fmt.Errorf("unsupported runtime: %s", runtime)
	}

	return cmd, nil
}

func initGoProject(path, projectName string) (*exec.Cmd, error) {
	modulePrompt := promptui.Prompt{
		Label: "Enter Go module name (e.g., github.com/username/project)",
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("module name cannot be empty")
			}
			return nil
		},
	}

	moduleName, err := modulePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return nil, err
	}

	err = createGoProjectStructure(path, moduleName)
	if err != nil {
		fmt.Printf("Error creating Go project structure: %v\n", err)
		return nil, err
	}

	return exec.Command("go", "mod", "init", moduleName), err
}

func createGoProjectStructure(projectPath, moduleName string) error {
	folders := []string{
		"cmd",
		"internal",
		"pkg",
		"api",
		"web",
		"configs",
		"deployments",
		"test",
		"docs",
		"tools",
		"scripts",
	}

	for _, folder := range folders {
		err := os.MkdirAll(filepath.Join(projectPath, folder), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating folder %s: %v", folder, err)
		}
	}

	// Create main.go in cmd/projectname
	projectName := filepath.Base(projectPath)
	cmdProjectDir := filepath.Join(projectPath, "cmd", projectName)
	err := os.MkdirAll(cmdProjectDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating cmd/%s directory: %v", projectName, err)
	}

	mainContent := fmt.Sprintf(`package main

import (
	"fmt"
	"%s/internal/app"
	"%s/internal/config"
	"%s/pkg/database"
	"%s/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)
	db := database.New(cfg.DatabaseURL)

	app := app.New(cfg, log, db)
	if err := app.Run(); err != nil {
		log.Error("Failed to run app", "error", err)
	}
}
`, moduleName, moduleName, moduleName, moduleName)

	err = os.WriteFile(filepath.Join(cmdProjectDir, "main.go"), []byte(mainContent), 0644)
	if err != nil {
		return fmt.Errorf("error creating main.go: %v", err)
	}

	// Create other necessary files
	files := map[string]string{
		filepath.Join("internal", "app", "app.go"): fmt.Sprintf(`package app

import (
	"%s/internal/config"
	"%s/pkg/database"
	"%s/pkg/logger"
)

type App struct {
	cfg *config.Config
	log logger.Logger
	db  *database.Database
}

func New(cfg *config.Config, log logger.Logger, db *database.Database) *App {
	return &App{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (a *App) Run() error {
	a.log.Info("Starting the application")
	// Add your application logic here
	return nil
}
`, moduleName, moduleName, moduleName),

		filepath.Join("internal", "config", "config.go"): `package config

type Config struct {
	LogLevel    string
	DatabaseURL string
}

func Load() *Config {
	// TODO: Implement config loading logic (e.g., from env vars or config file)
	return &Config{
		LogLevel:    "info",
		DatabaseURL: "postgres://user:password@localhost:5432/dbname",
	}
}
`,

		filepath.Join("pkg", "database", "database.go"): `package database

type Database struct {
	// Add database-specific fields here
}

func New(url string) *Database {
	// TODO: Implement database connection logic
	return &Database{}
}
`,

		filepath.Join("pkg", "logger", "logger.go"): `package logger

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

func New(level string) Logger {
	// TODO: Implement logger initialization logic
	return &defaultLogger{}
}

type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, keysAndValues ...interface{}) {
	// TODO: Implement info logging
}

func (l *defaultLogger) Error(msg string, keysAndValues ...interface{}) {
	// TODO: Implement error logging
}
`,
	}

	for path, content := range files {
		err := os.MkdirAll(filepath.Dir(filepath.Join(projectPath, path)), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory for %s: %v", path, err)
		}
		err = os.WriteFile(filepath.Join(projectPath, path), []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("error creating file %s: %v", path, err)
		}
	}

	return nil
}

func isExecutableAvailable(name string) bool {
	cmd := exec.Command(name, "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
