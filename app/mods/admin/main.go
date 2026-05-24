package admin

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"mthan-go-starter/app/routes"
	"mthan-go-starter/app/services"
)

// Register mounts the admin panel route with basic authentication.
func Register(mux *http.ServeMux, svc *services.ProcessService, logger *services.LoggerService, cfg *services.Config) {
	adminPath := cfg.AdminPath
	// Ensure it starts with /
	if adminPath[0] != '/' {
		adminPath = "/" + adminPath
	}

	logger.Info("Activating Admin Panel", "path", adminPath)

	adminBuildPath := filepath.Join("app", "mods", "admin", "build")

	// Ensure path ends with / for sub-route serving
	regPath := adminPath
	if !strings.HasSuffix(regPath, "/") {
		regPath = regPath + "/"
	}

	mux.HandleFunc(regPath, routes.LoggingMiddleware(logger, routes.RecoveryMiddleware(logger, routes.CORSMiddleware(
		routes.BasicAuth(cfg.AdminUser, cfg.AdminPass, func(w http.ResponseWriter, r *http.Request) {
			// Check if built admin client exists
			if _, err := os.Stat(adminBuildPath); os.IsNotExist(err) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Admin Portal Pending</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #0f172a; color: #f8fafc; padding: 40px; display: flex; align-items: center; justify-content: center; min-height: 80vh; }
        .card { background: #1e293b; padding: 32px; border-radius: 8px; max-width: 600px; text-align: center; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); border: 1px solid #334155; }
        h1 { color: #38bdf8; margin-top: 0; font-size: 1.8em; }
        code { background: #0f172a; padding: 4px 8px; border-radius: 4px; color: #f43f5e; font-family: monospace; font-size: 1.1em; }
        p { color: #94a3b8; line-height: 1.6; }
    </style>
</head>
<body>
    <div class="card">
        <h1>Admin Portal Build Pending</h1>
        <p>The admin panel folder exists, but the production build cannot be found.</p>
        <p>Please compile the admin React client by executing:</p>
        <p><code>./scripts/client build-admin</code></p>
    </div>
</body>
</html>
				`)
				return
			}

			// Clean relative sub-route path
			relPath := strings.TrimPrefix(r.URL.Path, adminPath)
			if relPath == "" || relPath == "/" {
				relPath = "index.html"
			}

			// Serve the requested static asset or fall back to index.html for SPA routing
			targetFile := filepath.Join(adminBuildPath, relPath)
			if _, err := os.Stat(targetFile); os.IsNotExist(err) {
				http.ServeFile(w, r, filepath.Join(adminBuildPath, "index.html"))
				return
			}

			http.ServeFile(w, r, targetFile)
		}),
	))))
}
