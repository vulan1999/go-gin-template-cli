package generator

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

//go:embed assets/*
var templateFiles embed.FS

type templateData struct {
	ModulePath string
}

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
		filepath.Join(projectName, "cmd", "api", "main.go"):              "assets/main.tmpl",
		filepath.Join(projectName, "internal", "routes", "example.go"):   "assets/routes.tmpl",
		filepath.Join(projectName, "internal", "handlers", "example.go"): "assets/handler.tmpl",
		filepath.Join(projectName, "config", "enviroment.go"):            "assets/config.tmpl",
	}

	data := templateData{
		ModulePath: modulePath,
	}

	for outputPath, templatePath := range files {
		if err := renderFromTemplate(outputPath, templatePath, data); err != nil {
			return err
		}
	}
	return nil
}

// renderFromTemplate write files base on tmpl on assets forlder
func renderFromTemplate(outputPath, templatePath string, data templateData) error {
	tmpl, err := template.ParseFS(templateFiles, templatePath)
	if err != nil {
		return fmt.Errorf("cannot parse template %s: %w", templatePath, err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create file %s: %w", outputPath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("cannot render data to template %s: %w", templatePath, err)
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
