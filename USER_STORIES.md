# User Stories — Go Gin Template CLI

> Format: each story states a role, a goal, and a motivation, followed by
> acceptance criteria in Given/When/Then form. This format is intentionally
> structured so both humans and AI tools can parse "what counts as done."
>
> Status legend: ✅ Implemented · 🚧 Planned

---

## US-1: Generate a barebone local project ✅

**As a** Go developer
**I want** to generate a new Gin project without a remote repo association
**So that** I can quickly start prototyping locally before deciding where to host the code

**Acceptance Criteria**
```gherkin
Given I run the CLI
When I enter a project name and select "Barebone" as the module convention
Then a new folder with that project name is created
And the go.mod module path equals the project name
And Gin and godotenv are installed as dependencies
And cmd/api/main.go, internal/routes/example.go, internal/handlers/example.go,
  and config/enviroment.go are created
```

---

## US-2: Generate a project tied to a GitHub repo ✅

**As a** Go developer
**I want** to generate a project whose module path matches my GitHub repo
**So that** `go get` and imports work correctly once I push the code to GitHub

**Acceptance Criteria**
```gherkin
Given I run the CLI and enter a project name
When I select "GitHub" as the module convention
And I enter my GitHub username or organization
Then the go.mod module path is set to github.com/<username>/<project-name>
And the rest of the project is generated the same as the Barebone flow
```

---

## US-3: Generate a project tied to a GitLab repo (including self-hosted) ✅

**As a** Go developer using GitLab (including a self-hosted instance)
**I want** to specify a custom GitLab host and my namespace
**So that** the generated module path is correct for my GitLab setup

**Acceptance Criteria**
```gherkin
Given I run the CLI and enter a project name
When I select "GitLab" as the module convention
Then I am prompted for a GitLab host, defaulting to "gitlab.com" if left blank
And I am then prompted for my username or group
And the go.mod module path is set to <host>/<username>/<project-name>
```

---

## US-4: Cancel generation at any point ✅

**As a** user of the CLI
**I want** to quit the wizard at any step without side effects
**So that** I don't end up with a half-created project if I change my mind

**Acceptance Criteria**
```gherkin
Given I am at any step of the wizard (name entry, module choice, or generating)
When I send a quit signal (e.g. Ctrl+C or Esc, per the input component's behavior)
Then the CLI exits
And no project files are left in a broken/partial state
```
*(Note: verify actual partial-state handling during generation — this is an
assumption to confirm against current `spinner` step behavior.)*

---

## US-5: See generation progress ✅

**As a** user running the generator
**I want** to see each generation step and whether it succeeded
**So that** I know what's happening and can diagnose failures (e.g. `go get` network issues)

**Acceptance Criteria**
```gherkin
Given I have completed the name and module-choice steps
When generation starts
Then I see a spinner with a live-updating title for each step:
  "Creating project directories...", "Initializing Go module...",
  "Installing Gin and dependencies...", "Creating template files..."
And each step's title updates to its "done" state on success
And if any step fails, the error is surfaced to me
```

---

## US-6: Generate a project with a database connection 🚧

**As a** Go developer
**I want** to choose a database (e.g. PostgreSQL, MySQL, SQLite) and have the
connection code generated for me
**So that** I don't have to hand-write config, connection pooling, and a
sample repository every time

**Acceptance Criteria (proposed)**
```gherkin
Given I am in the module choice / setup flow
When I am prompted to choose a database ("None", "PostgreSQL", "MySQL", "SQLite")
And I select a database other than "None"
Then a config/database.go file is generated with a connection function
  reading DSN values from environment variables
And an internal/repositories package is generated with one example
  repository demonstrating basic CRUD
And the appropriate driver/ORM dependency is installed via go get
And cmd/api/main.go is updated to initialize the DB connection on startup
```

---

## US-7: Get a working local database via Docker Compose 🚧

**As a** developer who just generated a project with a database
**I want** a docker-compose.yml for the chosen database
**So that** I can run `docker compose up` and immediately connect, without
installing a database locally

**Acceptance Criteria (proposed)**
```gherkin
Given I selected a database other than "None" during generation
When generation completes
Then a docker-compose.yml is created in the project root
  defining a service for the selected database
And the .env / config defaults match the docker-compose service's
  credentials and port
```

---

## How to Extend This Document

When adding a new feature, add a new `US-N` entry with:
1. Role / goal / motivation (the "As a / I want / So that" line)
2. Acceptance criteria in Given/When/Then form
3. Status marker (✅ / 🚧)

Keep each Given/When/Then block scoped to one behavior — split into multiple
stories rather than writing one story with many unrelated "When" clauses.

