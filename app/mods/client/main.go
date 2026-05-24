package client

import (
	"net/http"
	"path/filepath"

	"mthan-go-starter/app/services"
)

// Register mounts the public React client routing fallback.
func Register(mux *http.ServeMux, svc *services.ProcessService, logger *services.LoggerService, cfg *services.Config) {
	logger.Info("Serving client static files", "path", cfg.ClientBuildPath)
	fs := http.Dir(cfg.ClientBuildPath)
	fileServer := http.FileServer(fs)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(r.URL.Path)
		if err != nil {
			// File does not exist, serve index.html for SPA routing fallback
			http.ServeFile(w, r, filepath.Join(cfg.ClientBuildPath, "index.html"))
			return
		}
		f.Close()
		fileServer.ServeHTTP(w, r)
	})
}
