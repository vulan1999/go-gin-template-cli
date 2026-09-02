package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateDirectories creates the initial folder structure for the Gin project.
func CreateDirectories(projectName string) error {
	dirs := []string{
		filepath.Join(projectName, "cmd", "api"),
		filepath.Join(projectName, "internal", "routes"),
		filepath.Join(projectName, "internal", "handlers"),
		filepath.Join(projectName, "config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("fail to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// InitGoModule initializes the go module with the specified module path.
func InitGoModule(projectName, modulePath string) error {
	cmd := exec.Command("go", "mod", "init", modulePath)
	cmd.Dir = projectName
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cannot run go mod init: %s (%w)", string(out), err)
	}
	return nil
}

// InstallDependencies installs gin and godotenv packages.
func InstallDependencies(projectName string) error {
	getGinCmd := exec.Command("go", "get", "github.com/gin-gonic/gin")
	getGinCmd.Dir = projectName
	if out, err := getGinCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cannot get Gin: %s (%w)", string(out), err)
	}

	getGodotenvCmd := exec.Command("go", "get", "github.com/joho/godotenv")
	getGodotenvCmd.Dir = projectName
	if out, err := getGodotenvCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cannot get godotenv: %s (%w)", string(out), err)
	}
	return nil
}

// CreateTemplateFiles writes the boilerplate templates to the project directory.
func CreateTemplateFiles(projectName, modulePath string) error {
	files := map[string]string{
		filepath.Join(projectName, "cmd", "api", "main.go"):              getMainFileTemplate(modulePath),
		filepath.Join(projectName, "internal", "routes", "example.go"):   getRouteExampleTemplate(modulePath),
		filepath.Join(projectName, "internal", "handlers", "example.go"): getHandlerExampleTemplate(),
		filepath.Join(projectName, "config", "enviroment.go"):            getConfigFileTemplate(),
		filepath.Join(projectName, "Makefile"):                           getMakefileTemplate(),
	}

	for filePath, content := range files {
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("cannot write file %s: %w", filePath, err)
		}
	}
	return nil
}

// GenerateGinProject runs all generation steps sequentially.
func GenerateGinProject(projectName string, modulePath string) error {
	if err := CreateDirectories(projectName); err != nil {
		return err
	}
	if err := InitGoModule(projectName, modulePath); err != nil {
		return err
	}
	if err := InstallDependencies(projectName); err != nil {
		return err
	}
	if err := CreateTemplateFiles(projectName, modulePath); err != nil {
		return err
	}
	return nil
}
