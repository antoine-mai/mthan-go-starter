# Mthan Go Starter Kit

A high-performance Golang backend starter template integrated with a public React SPA client, an isolated Admin Control Panel React client, and a Django/Jinja2-like Pongo2 SSR page templating engine.

---

## 📁 Directory & Folder Structure

```text
mthan-go-starter/
├── app/                     # Backend Source Code
│   ├── mods/                # Self-contained modules (or plugins)
│   │   ├── admin/           # Admin Panel Module (Go handler + React Client)
│   │   └── client/          # Public React Client Module (Go handler + React Client)
│   ├── routes/              # Public HTML pages, API endpoints, and Action controllers
│   │   ├── api/             # Standard GET/POST API endpoints (e.g., /api/hello)
│   │   ├── pages/           # Server-Side Rendered (SSR) HTML pages (using templates)
│   │   ├── post/            # Action/Post same-origin controllers
│   │   └── main.go          # Common HTTP middlewares (CORS, Logging, Recovery)
│   └── services/            # Core business logic, configuration loader, and shared services
├── scripts/                 # Developer Utility Bash Scripts
│   ├── app                  # Unified app runner supporting {dev, run, test, build}
│   ├── client               # Unified client compiler supporting {build-app, build-admin, dev-app, dev-admin}
│   └── push                 # Utility script to stage, commit, and push modifications
├── templates/               # HTML templates for server pages
│   ├── pages/               # Individual content templates (e.g. home.html, contact.html)
│   └── base.html            # Global shared layout block base template
├── config.yaml              # Local development configuration file
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

## 🧩 Architectural Standards & Rules

### 1. Separation of Concerns
* **`routes/` Layer**: Responsible **only** for HTTP routing, request deserialization, payload validation, and response serialization. No business processing or domain logic is allowed here. Located in `app/routes/`.
* **`services/` Layer**: Responsible for all business processes, operations, database validation, and data logic. Handlers in the `routes/` layer delegate execution to this layer. Logging is also implemented as a service (`LoggerService` inside the `app/services/` layer) to allow clean dependency injection.
* **`mods/` Layer**: Houses self-contained plug-and-play features or larger optional application parts (like React clients, admin panel handlers, etc.).

### 2. Routing Conventions & Rules
* **`/api/...`**: Used for public backend APIs (defined in `app/routes/api/`).
* **`/post/...`**: Used for same-origin action/POST endpoints (defined in `app/routes/post/`).
* **SSR Pages**: Dynamic server-side page endpoints are defined in `app/routes/pages/` and rendered using template files located in `templates/pages/` inheriting from `templates/base.html` using Pongo2.
* **Fallback Routing Rule**:
  * If `client: true` is configured, the server serves the compiled React public client SPA assets from `app/mods/client/build/` at the root route `/` (with SPA fallback routing).
  * If `client: false` and `templates/` folder is present, the server renders Pongo2 HTML templates for page routes.
  * If templates are missing and client is false, the server falls back to mounting a default landing page handler from `app/mods/main.go`.

### 3. Avoiding Circular Dependencies
* The base `routes` package (`app/routes/main.go`) contains shared HTTP utilities (JSON response helpers, logging/recovery middlewares) and does **not** import any sub-routes.
* Sub-route packages (e.g. `app/routes/api/hello`, `app/routes/pages/contact`) import the base `routes` package to use the shared middlewares and helpers.
* The application entrypoint (`main.go`) imports all sub-routers and registers them sequentially on the HTTP multiplexer.

---

## 🧩 Module Development Guidelines

When developing modules in `app/mods/`, you can structure them in two ways depending on their purpose and lifecycle:

### 1. Self-Contained Modules (Recommended for Plug-and-Play features)
* **Design**: The module defines its own handlers and registers its routes directly.
* **Benefits**: **Complete isolation**. The module can be deleted at any time, and the code-generation compilation script (`./scripts/app build` or `./scripts/app dev`) will automatically discover the change and compile cleanly without leaving broken references in the core routes.
* **Example**: The `admin` module is a self-contained module.

### 2. Process/Service Modules (Recommended for Core business logic)
* **Design**: The module contains purely business logic (e.g., database transactions, calculations) and exposes functions. The core routes in `app/routes/` import the module and invoke these functions to handle API requests.
* **Benefits**: Centralized control of all API endpoints under the `app/routes/` directory.
* **Trade-off**: Deleting the module will require manually removing its import statements and calls from `app/routes/` to avoid compilation errors.

---

## 🛠️ Developer Workflows & Commands

Manage building and running modules using the provided scripts:

* **Concurrently Run Backend & Public Client in Dev Mode**:
  ```bash
  ./scripts/app dev
  ```
* **Build / Run React Clients**:
  * Build Public Client: `./scripts/client build-app` (Compiles assets to `app/mods/client/build/`)
  * Build Admin Client: `./scripts/client build-admin` (Compiles assets to `app/mods/admin/build/`)
  * Dev Public Client: `./scripts/client dev-app` (Vite dev server on port `3000`)
  * Dev Admin Client: `./scripts/client dev-admin` (Vite dev server on port `3001`)
* **Compile Go Backend**:
  ```bash
  ./scripts/app build
  ```
* **Run Go Backend Binary**:
  ```bash
  ./scripts/app run
  ```
  Pass flags to the executable directly. For example, to run in standalone mode:
  ```bash
  ./scripts/app run --standalone
  ```
* **Run Go Tests**:
  ```bash
  ./scripts/app test
  ```