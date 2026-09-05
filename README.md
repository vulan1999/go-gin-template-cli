# Go Gin Template Generate CLI

## 1. What This Project Is

A command-line tool that scaffolds a new [Gin](https://github.com/gin-gonic/gin)
web API project in Go. Instead of manually creating folders, writing
boilerplate handlers, and wiring up `go.mod`, a developer runs the CLI,
answers a few interactive prompts, and gets a working project skeleton.

The CLI's interactive prompts are built with
[bubbletea](https://github.com/charmbracelet/bubbletea) (a terminal UI
framework), and the actual file/folder generation logic is plain Go using
`text/template` and embedded template files.

## 2. Goals

- Save developers the repetitive setup work of starting a new Gin API.
- Enforce a consistent, opinionated project layout across projects.
- Keep the generated code minimal and idiomatic — no hidden magic, easy to
  read and extend by hand after generation.

## 3. Non-Goals (for now)

- This is not a full backend framework — it generates a *starting point*,
  not a batteries-included framework the user imports as a dependency.
- Not currently opinionated about deployment, CI/CD, or cloud provider.

## 4. How It Works (User Flow)

1. User runs the CLI binary.
2. **Step 1 — Project name.** A text prompt asks for the project name
   (used as the output folder name).
3. **Step 2 — Module naming convention.** User picks one of:
   - **Barebone** — module path = project name only (local-only module).
   - **GitHub** — user enters a GitHub username/org; module path becomes
     `github.com/<user>/<project>`.
   - **GitLab** — user enters a GitLab host (defaults to `gitlab.com`) and a
     username/group; module path becomes `<host>/<user>/<project>`.
4. **Step 3 — Generation.** A spinner runs through these steps in order:
   1. Create project directories
   2. `go mod init <modulePath>`
   3. `go get` Gin and godotenv
   4. Render template files into the new project
5. Done — user has a working Gin project on disk.

## 5. Generated Project Structure

```text
<project-name>/
├── cmd/
│   └── api/
│       └── main.go            # Entry point
├── config/
│   └── enviroment.go          # Reads environment variables (godotenv)
├── internal/
│   ├── routes/
│   │   └── example.go         # Route definitions
│   └── handlers/
│       └── example.go         # HTTP request/response handling
├── go.mod
└── go.sum
```

> Note: `internal/models`, `internal/repositories`, and `internal/services`
> are documented as intended structure but are **not yet generated** by the
> CLI as of the current version — see Roadmap.

## 6. Source Repo Structure

```text
go-gin-template-cli/
├── main.go                    # CLI entry point
├── generator/
│   ├── generator.go           # Core generation logic (dirs, go mod, deps, templates)
│   └── assets/                # .tmpl files embedded via go:embed
├── ui/
│   ├── wizard.go               # Top-level bubbletea flow (name -> module -> generate)
│   ├── module_choice.go        # Module naming convention sub-flow
│   └── components/             # Reusable TUI pieces (text input, choice list, spinner)
├── Makefile
├── go.mod / go.sum
└── README.md
```

## 7. Current Features (Implemented)

- Interactive project name prompt
- Module path selection: Barebone / GitHub / GitLab
- Automatic `go mod init` with the chosen module path
- Automatic install of `gin-gonic/gin` and `joho/godotenv`
- Generated files: `cmd/api/main.go`, `internal/routes/example.go`,
  `internal/handlers/example.go`, `config/enviroment.go`

## 8. Roadmap / Not Yet Implemented

- [x] Database connection setup (driver + ORM choice, e.g. Postgres/MySQL/SQLite
      with GORM or `database/sql`)
- [x] Generated `internal/models` and `internal/repositories` layers
- [x] Optional Docker Compose file for local database
- [x] Service layer scaffolding (`internal/services`)

## 9. Tech Stack

| Concern            | Tool/Library                          |
|---------------------|----------------------------------------|
| Language            | Go                                     |
| Web framework (generated projects) | Gin (`gin-gonic/gin`)   |
| Env var loading (generated projects) | godotenv (`joho/godotenv`) |
| CLI terminal UI      | bubbletea (`charmbracelet/bubbletea`) |
| Templating           | Go `text/template` + `embed`          |

## 10. Glossary

- **Module path** — the Go module identity (in `go.mod`), also used as the
  import path prefix for the generated project.
- **Barebone** — module naming mode with no remote host prefix.
- **Template (.tmpl)** — a file in `generator/assets` rendered with
  `templateData` (currently just `ModulePath`) to produce a real Go file.
