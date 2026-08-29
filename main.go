package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	// Get project name
	fmt.Print("📦 Enter your project name: ")
	projectNameInput, _ := reader.ReadString('\n')
	projectName := strings.TrimSpace(projectNameInput)

	if projectName == "" {
		fmt.Println("Error: Project name cannot be empty.")
		os.Exit(1)
	}

	// Initialize go mod
	fmt.Println("\nSelect your module naming convention:")
	fmt.Println("	[1] Barebone (Local only -> e.g., myapp)")
	fmt.Println("	[2] Github (Remote Github repo -> e.g., github.com/username/myapp)")
	fmt.Println("	[3] Gitlab (Remote Gitlab repo -> e.g., gitlab.com/username/myapp)")
	fmt.Print("Choose an option [1-3] (default 1): ")

	choiceInput, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(choiceInput)

	var modulePath string

	switch choice {
	// For github
	case "2":
		platform := "github.com"
		fmt.Printf("Enter your %s username or organiztion: ", platform)
		usernameInput, _ := reader.ReadString('\n')
		username := strings.TrimSpace(usernameInput)
		if username == "" {
			fmt.Println("Error: Username cannot be empty for remote repository")
			os.Exit(1)
		}

		modulePath = fmt.Sprintf("github.com/%s/%s", username, projectName)
	// For gitlab
	case "3":
		fmt.Print("Enter Gitlab domain (leave empty for gitlab.com): ")
		domainInput, _ := reader.ReadString('\n')
		gitlabDomain := strings.TrimSpace(domainInput)
		if gitlabDomain == "" {
			gitlabDomain = "gitlab.com"
		}

		fmt.Printf("Enter your username or organiztion for %s: ", gitlabDomain)
		usernameInput, _ := reader.ReadString('\n')
		username := strings.TrimSpace(usernameInput)

		if username == "" {
			fmt.Println("Error: Username cannot be empty for remote repository")
			os.Exit(1)
		}

		modulePath = fmt.Sprintf("%s/%s/%s", gitlabDomain, username, projectName)
	// Default choice
	default:
		modulePath = projectName
	}

	fmt.Printf("🚀 Generating Gin project: %s...\n", projectName)
	fmt.Printf("📁 Local Folder: %s\n", projectName)
	fmt.Printf("🏷️  Module Path:  %s\n\n", modulePath)

	// Create project root directory
	dirs := []string{
		filepath.Join(projectName, "internal", "handlers"),
		filepath.Join(projectName, "internal", "routes"),
		filepath.Join(projectName, "cmd", "server"),
	}

	// Create folder with following permission drwxr-xr-x
	// Owner: Full permission (read - write - execute)(7)
	// Group and Other user using this machine (read - execute)(5)
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Initialize go mod
	cmd := exec.Command("go", "mod", "init", modulePath)
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to run go mod init: %v\n", err)
		os.Exit(1)
	}

	// Initialize git
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = projectName
	if err := gitCmd.Run(); err != nil {
		fmt.Printf("Failed to run git init: %v", err)
		os.Exit(1)
	}

	// Generate boilerplate files
}
