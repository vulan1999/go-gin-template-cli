package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GenerateGinProject(projectName string, modulePath string) error {
	// Make Directory Tree
	dirs := []string{
		filepath.Join(projectName, "cmd", "api"),
		filepath.Join(projectName, "internal", "routes"),
		filepath.Join(projectName, "internal", "services"),
		filepath.Join(projectName, "config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("fail to create directory %s: %v", dir, err)
		}
	}
	// Run go mod init command with the module name
	cmd := exec.Command("go", "mod", "init", modulePath)
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cannot run go mod init: %v", err)
	}
	// Install Gin
	getGinCmd := exec.Command("go", "get", "github.com/gin-gonic/gin")
	getGinCmd.Dir = projectName
	if err := getGinCmd.Run(); err != nil {
		return fmt.Errorf("cannot get Gin: %v", err)
	}
	// Make Standard file template (main.go, README.md, healthcheck.go, etc...)
	files := map[string]string{
		filepath.Join(projectName, "cmd", "api", "main.go"):              "content",
		filepath.Join(projectName, "internal", "routes", "example.go"):   "content",
		filepath.Join(projectName, "internal", "services", "example.go"): "contnet",
	}

	for filepath, content := range files {
		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			return fmt.Errorf("cannot write file %s: %v", filepath, err)
		}
	}
	return nil
}
