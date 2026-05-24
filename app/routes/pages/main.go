package pages

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/flosch/pongo2/v6"

	"mthan-go-starter/app/services"
)

// Register mounts the HTML page routes.
func Register(mux *http.ServeMux, svc *services.ProcessService, logger *services.LoggerService, cfg *services.Config) {
	logger.Info("Registering SSR HTML Page routes")

	fs := http.Dir("templates")
	fileServer := http.FileServer(fs)

	// Homepage / and static assets handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If requesting a static sub-path, serve it from templates/ directly
		if r.URL.Path != "/" {
			f, err := fs.Open(r.URL.Path)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// Detect active modules for display
		var activeMods []string
		if _, err := os.Stat(filepath.Join("app", "mods", "admin")); err == nil {
			activeMods = append(activeMods, "Admin Panel (React)")
		}
		if _, err := os.Stat(filepath.Join("app", "mods", "client")); err == nil {
			activeMods = append(activeMods, "Public Client (React)")
		}

		ctx := pongo2.Context{
			"AppName":    cfg.AppName,
			"Version":    "v1.0.0",
			"Time":       time.Now().Format("15:04:05 PM"),
			"AdminPath":  cfg.AdminPath,
			"ActiveMods": activeMods,
		}

		renderTemplate(w, logger, "pages/home.html", ctx)
	})
}

// renderTemplate compiles the pongo2 template and executes it with context.
func renderTemplate(w http.ResponseWriter, logger *services.LoggerService, page string, ctx pongo2.Context) {
	pagePath := filepath.Join("templates", page)

	tmpl, err := pongo2.FromFile(pagePath)
	if err != nil {
		logger.Error("Failed to parse Pongo2 template", "error", err, "page", page)
		http.Error(w, "Template Parsing Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteWriter(ctx, w)
	if err != nil {
		logger.Error("Failed to execute Pongo2 template", "error", err, "page", page)
	}
}
