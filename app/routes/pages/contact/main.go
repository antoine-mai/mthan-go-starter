package contact

import (
	"net/http"
	"path/filepath"

	"github.com/flosch/pongo2/v6"

	"mthan-go-starter/app/services"
)

// Register mounts the /contact route on the ServeMux.
func Register(mux *http.ServeMux, svc *services.ProcessService, logger *services.LoggerService, cfg *services.Config) {
	mux.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		pagePath := filepath.Join("templates", "pages", "contact.html")

		tmpl, err := pongo2.FromFile(pagePath)
		if err != nil {
			logger.Error("Failed to parse support Pongo2 templates", "error", err)
			http.Error(w, "Template Parsing Error", http.StatusInternalServerError)
			return
		}

		ctx := pongo2.Context{
			"AppName": cfg.AppName,
			"Version": "v1.0.0",
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.ExecuteWriter(ctx, w)
		if err != nil {
			logger.Error("Failed to execute support Pongo2 template", "error", err)
		}
	})
}
