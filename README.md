# Mthan Go Starter Kit

A high-performance Golang backend starter template integrated with a public React SPA client and a dedicated, isolated Admin Control Panel React client.

---

## 📁 Directory & Folder Structure

```text
mthan-go-starter/
├── app/                     # Backend Source Code
│   ├── config/              # Configuration module (YAML parser, path resolver)
│   ├── mods/                # Modules
│   ├── routes/              # Public API / Action endpoint controllers
│   │   ├── api/             # Standard GET/POST API endpoints (e.g., /api/hello)
│   │   ├── post/            # Action/Post same-origin controllers
│   │   └── main.go          # Common HTTP middlewares (CORS, Logging, Recovery)
│   └── services/            # Core business logic layer and shared services
├── client/                  # Main Public React SPA Client
│   ├── src/                 # Client React source code
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── scripts/                 # Developer Utility Bash Scripts
│   ├── app                  # Unified developer tool supporting {dev, run, test, build}
│   ├── client               # Build / Run developer tasks for the Public client
│   └── admin                # Build / Run developer tasks for the Admin client
├── config.yaml              # Local development configuration template
├── go.mod                   # Go module definitions
├── main.go                  # Main entry point bootstrapping the application
└── README.md                # Project documentation
```

---

## ⚙️ Configuration (`config.yaml`)

The project uses a structured `config.yaml` file to configure all parameters:

```yaml
app:
  name: "mthan-app"          # App name, used to define Linux User config paths

server:
  port: "8080"               # HTTP port for the backend server
  env: "development"         # Runtime mode (development / production)

client: false                # Set true to serve the built public React client

database:
  driver: "sqlite"           # Database driver ("sqlite" or "postgres")
  url: ""                    # Connection URL (empty defaults SQLite to storage/db.sqlite)

admin:
  active: true               # Set true to activate the Admin Portal
  path: "/admin"             # Mount path for the admin portal
  user: "admin"              # Admin login username
  pass: "admin@123"          # Admin login password
```

---

## 🚀 Runtime Environments & Path Resolutions

The application automatically resolves path destinations depending on the running environment:

1. **Local Development Mode** (automatically detected when `go.mod` is present in the current working directory):
   * Configuration is read directly from `./config.yaml` at the root project folder.
   * Data storage directories are mapped to `./storage/`.
2. **Standalone Mode** (triggered by launching the compiled binary with `-standalone` or `--standalone` flags):
   * Configuration is read from `config/config.yaml` next to the executable.
   * Data storage is resolved next to the executable at `storage/`.
3. **Linux User Mode** (default fallback on Linux systems):
   * Configuration is read from `~/.[app.name]/config/config.yaml`.
   * Data storage is resolved to `~/.[app.name]/storage/`.

---

## 🛠️ Developer Workflows

Manage building and running modules using the provided scripts:

* **Concurrently Run Backend & Public Client in Dev Mode**:
  ```bash
  ./scripts/app dev
  ```
* **Build / Run Public React Client**:
  * Dev: `./scripts/client dev` (Vite dev server on port `3000`)
  * Build: `./scripts/client build` (Compiles assets into `client/build/`)
* **Build / Run Admin React Client**:
  * Dev: `./scripts/admin dev` (Vite dev server on port `3001`)
  * Build: `./scripts/admin build` (Compiles assets into `app/mods/admin/client/build/`)
* **Compile Go Backend**:
  ```bash
  ./scripts/app build
  ```
* **Run Go Backend Binary**:
  ```bash
  ./scripts/app run
  ```
* **Run Go Tests**:
  ```bash
  ./scripts/app test
  ```