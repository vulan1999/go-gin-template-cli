package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateBaseTemplateFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-gen-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projectName := filepath.Join(tempDir, "my-app")
	if err := CreateDirectories(projectName); err != nil {
		t.Fatalf("CreateDirectories failed: %v", err)
	}

	modulePath := "github.com/test/my-app"
	if err := CreateBaseTemplateFiles(projectName, modulePath, "PostgreSQL", "GORM"); err != nil {
		t.Fatalf("CreateBaseTemplateFiles failed: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(projectName, "cmd", "api", "main.go"),
		filepath.Join(projectName, "internal", "routes", "example.go"),
		filepath.Join(projectName, "internal", "handlers", "example.go"),
		filepath.Join(projectName, "config", "enviroment.go"),
		filepath.Join(projectName, ".env.example"),
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}

	// Verify main.go contains InitDB when database is set
	mainContent, err := os.ReadFile(filepath.Join(projectName, "cmd", "api", "main.go"))
	if err != nil {
		t.Fatalf("failed to read main.go: %v", err)
	}
	if !strings.Contains(string(mainContent), "config.InitDB()") {
		t.Errorf("expected main.go to contain config.InitDB(), got:\n%s", string(mainContent))
	}

	// Verify .env.example contains DB credentials
	envContent, err := os.ReadFile(filepath.Join(projectName, ".env.example"))
	if err != nil {
		t.Fatalf("failed to read .env.example: %v", err)
	}
	if !strings.Contains(string(envContent), "DB_HOST=localhost") {
		t.Errorf("expected .env.example to contain DB_HOST, got:\n%s", string(envContent))
	}
}

func TestCreateBaseTemplateFiles_NoneDatabase(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-gen-none-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projectName := filepath.Join(tempDir, "my-app")
	if err := CreateDirectories(projectName); err != nil {
		t.Fatalf("CreateDirectories failed: %v", err)
	}

	modulePath := "github.com/test/my-app"
	if err := CreateBaseTemplateFiles(projectName, modulePath, "None", ""); err != nil {
		t.Fatalf("CreateBaseTemplateFiles failed: %v", err)
	}

	mainContent, err := os.ReadFile(filepath.Join(projectName, "cmd", "api", "main.go"))
	if err != nil {
		t.Fatalf("failed to read main.go: %v", err)
	}
	if strings.Contains(string(mainContent), "config.InitDB()") {
		t.Errorf("expected main.go NOT to contain config.InitDB(), got:\n%s", string(mainContent))
	}
}

func TestCreateDatabaseTemplateFiles(t *testing.T) {
	tests := []struct {
		name          string
		database      string
		toolkit       string
		expectedInDB  []string
	}{
		{
			name:     "PostgreSQL + GORM",
			database: "PostgreSQL",
			toolkit:  "GORM",
			expectedInDB: []string{
				"gorm.io/gorm",
				"gorm.io/driver/postgres",
				"var DB *gorm.DB",
				"gorm.Open(postgres.Open(dsn)",
			},
		},
		{
			name:     "MySQL + GORM",
			database: "MySQL",
			toolkit:  "GORM",
			expectedInDB: []string{
				"gorm.io/gorm",
				"gorm.io/driver/mysql",
				"var DB *gorm.DB",
				"gorm.Open(mysql.Open(dsn)",
			},
		},
		{
			name:     "SQLite + GORM",
			database: "SQLite",
			toolkit:  "GORM",
			expectedInDB: []string{
				"gorm.io/gorm",
				"gorm.io/driver/sqlite",
				"var DB *gorm.DB",
				"gorm.Open(sqlite.Open(dbPath)",
			},
		},
		{
			name:     "PostgreSQL + sqlx",
			database: "PostgreSQL",
			toolkit:  "sqlx",
			expectedInDB: []string{
				"github.com/jmoiron/sqlx",
				"github.com/lib/pq",
				"var DB *sqlx.DB",
				"sqlx.Connect(driverName, dsn)",
			},
		},
		{
			name:     "MySQL + sqlx",
			database: "MySQL",
			toolkit:  "sqlx",
			expectedInDB: []string{
				"github.com/jmoiron/sqlx",
				"github.com/go-sql-driver/mysql",
				"var DB *sqlx.DB",
				"sqlx.Connect(driverName, dsn)",
			},
		},
		{
			name:     "SQLite + sqlx",
			database: "SQLite",
			toolkit:  "sqlx",
			expectedInDB: []string{
				"github.com/jmoiron/sqlx",
				"github.com/mattn/go-sqlite3",
				"var DB *sqlx.DB",
				"sqlx.Connect(driverName, dsn)",
			},
		},
		{
			name:     "PostgreSQL + database/sql",
			database: "PostgreSQL",
			toolkit:  "database/sql",
			expectedInDB: []string{
				"database/sql",
				"github.com/lib/pq",
				"var DB *sql.DB",
				"sql.Open(driverName, dsn)",
			},
		},
		{
			name:     "MySQL + database/sql",
			database: "MySQL",
			toolkit:  "database/sql",
			expectedInDB: []string{
				"database/sql",
				"github.com/go-sql-driver/mysql",
				"var DB *sql.DB",
				"sql.Open(driverName, dsn)",
			},
		},
		{
			name:     "SQLite + database/sql",
			database: "SQLite",
			toolkit:  "database/sql",
			expectedInDB: []string{
				"database/sql",
				"github.com/mattn/go-sqlite3",
				"var DB *sql.DB",
				"sql.Open(driverName, dsn)",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "test-db-tmpl-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			projectName := filepath.Join(tempDir, "my-app")
			if err := CreateDirectories(projectName); err != nil {
				t.Fatalf("CreateDirectories failed: %v", err)
			}

			err = CreateDatabaseTemplateFiles(projectName, "github.com/test/my-app", tc.database, tc.toolkit)
			if err != nil {
				t.Fatalf("CreateDatabaseTemplateFiles failed: %v", err)
			}

			dbFile := filepath.Join(projectName, "config", "database.go")
			content, err := os.ReadFile(dbFile)
			if err != nil {
				t.Fatalf("failed to read database.go: %v", err)
			}

			strContent := string(content)
			for _, exp := range tc.expectedInDB {
				if !strings.Contains(strContent, exp) {
					t.Errorf("expected database.go to contain %q, but got:\n%s", exp, strContent)
				}
			}
		})
	}
}

func TestCreateDatabaseTemplateFiles_None(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-db-none-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projectName := filepath.Join(tempDir, "my-app")
	if err := CreateDirectories(projectName); err != nil {
		t.Fatalf("CreateDirectories failed: %v", err)
	}

	if err := CreateDatabaseTemplateFiles(projectName, "github.com/test/my-app", "None", ""); err != nil {
		t.Fatalf("unexpected error for None: %v", err)
	}

	dbFile := filepath.Join(projectName, "config", "database.go")
	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		t.Errorf("database.go should not exist when None is chosen")
	}
}
