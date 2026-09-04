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
	Database   string
	Toolkit    string
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

func InitGitRepository(projectName string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = projectName
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cannot initialize git repository: %s (%w)", string(out), err)
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

// CreateBaseTemplateFiles writes the boilerplate base templates to the project directory. (main, route, handle, config, env files)
func CreateBaseTemplateFiles(projectName, modulePath, database, toolkit string) error {
	files := map[string]string{
		filepath.Join(projectName, "cmd", "api", "main.go"):              "assets/main.tmpl",
		filepath.Join(projectName, "internal", "routes", "example.go"):   "assets/routes.tmpl",
		filepath.Join(projectName, "internal", "handlers", "example.go"): "assets/handler.tmpl",
		filepath.Join(projectName, "config", "enviroment.go"):            "assets/config.tmpl",
		filepath.Join(projectName, ".env.example"):                       "assets/env.tmpl",
	}

	data := templateData{
		ModulePath: modulePath,
		Database:   database,
		Toolkit:    toolkit,
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

// InstallDatabaseDriverAndToolkit installs database drivers and toolkit packages based on user selection.
func InstallDatabaseDriverAndToolkit(projectName, database, toolkit string) error {
	if database == "" || database == "None" {
		return nil
	}

	var packages []string

	switch toolkit {
	case "GORM":
		packages = append(packages, "gorm.io/gorm")
		switch database {
		case "PostgreSQL":
			packages = append(packages, "gorm.io/driver/postgres")
		case "MySQL":
			packages = append(packages, "gorm.io/driver/mysql")
		case "SQLite":
			packages = append(packages, "gorm.io/driver/sqlite")
		}
	case "sqlx":
		packages = append(packages, "github.com/jmoiron/sqlx")
		switch database {
		case "PostgreSQL":
			packages = append(packages, "github.com/lib/pq")
		case "MySQL":
			packages = append(packages, "github.com/go-sql-driver/mysql")
		case "SQLite":
			packages = append(packages, "github.com/mattn/go-sqlite3")
		}
	case "database/sql":
		switch database {
		case "PostgreSQL":
			packages = append(packages, "github.com/lib/pq")
		case "MySQL":
			packages = append(packages, "github.com/go-sql-driver/mysql")
		case "SQLite":
			packages = append(packages, "github.com/mattn/go-sqlite3")
		}
	}

	for _, pkg := range packages {
		cmd := exec.Command("go", "get", pkg)
		cmd.Dir = projectName
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("cannot get %s: %s (%w)", pkg, string(out), err)
		}
	}

	return nil
}

// CreateDatabaseTemplateFiles generates database configuration and boilerplate files.
func CreateDatabaseTemplateFiles(projectName, modulePath, database, toolkit string) error {
	if database == "" || database == "None" {
		return nil
	}
	// TODO: generate repository file, model and service files
	var templatePath string
	switch toolkit {
	case "GORM":
		templatePath = "assets/database_gorm.tmpl"
	case "sqlx":
		templatePath = "assets/database_sqlx.tmpl"
	case "database/sql":
		templatePath = "assets/database_sql.tmpl"
	default:
		return fmt.Errorf("unsupported toolkit: %s", toolkit)
	}

	data := templateData{
		ModulePath: modulePath,
		Database:   database,
		Toolkit:    toolkit,
	}

	outputPath := filepath.Join(projectName, "config", "database.go")
	return renderFromTemplate(outputPath, templatePath, data)
}
