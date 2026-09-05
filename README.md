# Go Gin Template Generate CLI

A CLI tool that auto-generates a [Gin](https://github.com/gin-gonic/gin) web
API project for you. It uses [bubbletea](https://github.com/charmbracelet/bubbletea)
for an interactive terminal UI, so you just answer a few prompts and get a
ready-to-run Gin project — no manual folder setup, no hand-written boilerplate.
 
## Requirements
 
- [Go](https://go.dev/dl/) **1.24.2** or later
- Git (only needed if you plan to push the generated project to GitHub/GitLab)
- Internet access when generating a project (the CLI runs `go get` to install
  Gin and other dependencies into your new project)
## Installation
 
### Option 1: Install with `go install` (recommended)
 
```bash
go install github.com/vulan1999/go-gin-template-cli@latest
```
 
This places a `go-gin-template-cli` binary in your `$GOPATH/bin` (or
`$GOBIN`). Make sure that directory is on your `PATH`, then run it from
anywhere:
 
```bash
go-gin-template-cli
```
 
### Option 2: Clone and build locally
 
```bash
git clone https://github.com/vulan1999/go-gin-template-cli.git
cd go-gin-template-cli
make build
```
 
This produces a binary at `bin/go-gin-template-cli`. Run it with:
 
```bash
./bin/go-gin-template-cli
```
 
### Option 3: Run without building (for trying it out or development)
 
```bash
git clone https://github.com/vulan1999/go-gin-template-cli.git
cd go-gin-template-cli
make run
```
 
## Usage
 
Run the CLI, then follow the interactive prompts:
 
```bash
go-gin-template-cli
```
 
**Step 1 — Project name**
Enter the name of your project. This becomes the output folder name.
 
**Step 2 — Module naming convention**
Choose how your Go module path should be built:
 
| Option | Result | Example |
|---|---|---|
| **Barebone** | Local-only module, no remote prefix | `myproject` |
| **Github** | Prefixed with your GitHub username/org | `github.com/username/myproject` |
| **Gitlab** | Prefixed with your GitLab host + namespace | `gitlab.com/username/myproject` |
 
If you choose GitHub or GitLab, you'll be asked for your username/org (and,
for GitLab, the host — press Enter to default to `gitlab.com`).
 
**Step 3 — Generation**
The CLI will automatically:
1. Create the project directory structure
2. Run `go mod init <your-module-path>`
3. Install Gin and other required dependencies
4. Generate starter files from templates
You can cancel at any prompt with `Ctrl+C` or `Esc` — nothing is written
until generation starts.

## What You Get
 
Once generation finishes, your new project looks like this:
 
```text
<project-name>/
├── cmd/
│   └── api/
│       └── main.go            # Entry point — starts the Gin server
├── config/
│   └── enviroment.go          # Loads environment variables (.env support via godotenv)
├── internal/
│   ├── routes/
│   │   └── example.go         # Route definitions
│   └── handlers/
│       └── example.go         # Request/response handling
├── go.mod
└── go.sum
```
 
To run your generated project:
 
```bash
cd <project-name>
go run cmd/api/main.go
```


## License
 
See [LICENSE](./LICENSE).
